package common

import (
	"net"
	"net/http"
	"net/url"

	"strings"

	"net/http/cookiejar"
	"time"

	"context"

	"golang.org/x/net/proxy"
	"crawler-apkmirror/waiter"
	"github.com/golang/glog"
)

//HTTPClient 自带waiter限速
type HTTPClient struct {
	ProxyAddr string
	Client    http.Client
	Waiter    *waiter.BurstLimitTick
}

const (
	DefaultIdleTimeout    = 5 * time.Second
	DefaultConnectTimeout = 5 * time.Second
)

//TimeoutConn wraps a net.Conn, and sets a deadline for every read and write operation.
type TimeoutConn struct {
	net.Conn
	IdleTimeout time.Duration
}

/////////////////////////
//使用自定义出口协议,注意,前缀要全部使用小写
//如果是代理,那么使用 http:// 或者 https:// 类型的地址,如果使用出口 IP, 那么直接使用 ip:// 作为前缀
//如果使用ipv6, 那么使用`[]`把地址包起来
//例如:
//		http://14845132.xgj.me:27035
//		socks5://14845132.xgj.me:27035
//		ip://192.168.1.12
//		ip://[2607:5300:60:6566::]
func MakeTransportX(addr string) (transport *http.Transport) {
	transport = new(http.Transport)
	transport.MaxIdleConnsPerHost = 16
	//disable verify ssl
	//transport.TLSClientConfig = &tls.Config{
	//	InsecureSkipVerify: true,
	//}
	var (
		localAddr string
		dialer    proxy.Dialer
	)
	addr = strings.TrimSpace(addr)
	if strings.HasPrefix(addr, "ip") {
		localAddr = addr[5:]
	} else if strings.HasPrefix(addr, "http") {
		u := url.URL{}
		proxyUrl, e := u.Parse(addr)
		if e != nil {
			glog.V(0).Infof("set proxy failed, e:%v", e)
		} else {
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
	} else if strings.HasPrefix(addr, "socks5") {
		proxyUrl, err := url.Parse(addr)
		if err != nil {
			glog.V(0).Infof("Invalid proxy url %q", addr)
		}
		dialer, err = proxy.FromURL(proxyUrl, proxy.Direct)
		if err != nil {
			glog.V(0).Infof("proxy.FromURL %v err %v", proxyUrl, err)
		}
	} else if addr != "" {
		glog.V(0).Infof("MakeTransportX, addr (%s) have wrong format.", addr)
	}
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		var (
			conn net.Conn
		)
		if dialer != nil {
			var err error
			conn, err = dialer.Dial(network, addr)
			if err != nil {
				return nil, err
			}
		} else {
			d := net.Dialer{Timeout: DefaultConnectTimeout}
			if localAddr != "" && localAddr[0] == '[' {
				//如果本地ip地址以"["开头, 那么是ipv6地址，强制使用 tcp6 拨号
				network = "tcp6"
			}
			lAddr, err := net.ResolveTCPAddr(network, localAddr+":0")
			if err != nil {
				return nil, err
			}
			d.LocalAddr = lAddr
			conn, err = d.DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}
		}

		return NewTimeoutConn(conn, DefaultIdleTimeout)
	}
	return transport

}

//NewTimeoutConn 分配conn, 并指定超时
func NewTimeoutConn(conn net.Conn, idleTimeout time.Duration) (net.Conn, error) {
	c := &TimeoutConn{
		Conn:        conn,
		IdleTimeout: idleTimeout,
	}
	if c.IdleTimeout > 0 {
		deadline := time.Now().Add(idleTimeout)
		if e := c.Conn.SetDeadline(deadline); e != nil {
			return nil, e
		}
	}
	return c, nil
}
func (c *TimeoutConn) Read(b []byte) (int, error) {
	n, e := c.Conn.Read(b)
	if c.IdleTimeout > 0 && n > 0 && e == nil {
		err := c.Conn.SetDeadline(time.Now().Add(c.IdleTimeout))
		if err != nil {
			return 0, err
		}
	}
	return n, e
}

func (c *TimeoutConn) Write(b []byte) (int, error) {
	n, e := c.Conn.Write(b)
	if c.IdleTimeout > 0 && n > 0 && e == nil {
		err := c.Conn.SetDeadline(time.Now().Add(c.IdleTimeout))
		if err != nil {
			return 0, err
		}
	}
	return n, e
}

//NewHTTPClient 通过传入ip, 分配httpcli
func NewHTTPClient(proxyAddr string) *HTTPClient {
	c := &HTTPClient{ProxyAddr: proxyAddr}
	//Follow redirect 时复制 http header
	c.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		for attr, val := range via[0].Header {
			if _, ok := req.Header[attr]; !ok {
				req.Header[attr] = val
			}
		}
		return nil
	}

	//	c.Client.Timeout = 75 * time.Second
	return c
}

func (hc *HTTPClient) mkTransport() {
	if hc.Client.Transport != nil || hc.ProxyAddr == "" {
		return
	}
	hc.Client.Transport = MakeTransportX(hc.ProxyAddr)
}

//Do 完成出口ip &&发送http请求
func (hc *HTTPClient) Do(req *http.Request) (resp *http.Response, err error) {
	<-hc.Waiter.GetC()
	hc.mkTransport()
	return hc.Client.Do(req)
}

//EnableCookie 使能cookiee
func (hc *HTTPClient) EnableCookie() {
	if hc.Client.Jar == nil {
		cookieJar, _ := cookiejar.New(nil)
		hc.Client.Jar = cookieJar
	}
}

//DisableCookie 不使能cookie
func (hc *HTTPClient) DisableCookie() {
	hc.Client.Jar = nil
}

//IsCookieEnabled 是否使能cookie
func (hc *HTTPClient) IsCookieEnabled() bool {
	return hc.Client.Jar != nil
}

//GetCookies 获取cookies
func (hc *HTTPClient) GetCookies(u *url.URL) []*http.Cookie {
	if hc.Client.Jar == nil {
		return nil
	}
	return hc.Client.Jar.Cookies(u)
}

//GetCookie 获取cookie
func (hc *HTTPClient) GetCookie(u *url.URL, key string) *http.Cookie {
	if hc.Client.Jar == nil {
		return nil
	}
	//	u, e := url.Parse(rawUrl)
	//	if  e != nil {
	//		return nil
	//	}
	for _, c := range hc.Client.Jar.Cookies(u) {
		if c.Name == "key" {
			return c
		}
	}
	return nil
}
