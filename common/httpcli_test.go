package common

import (

	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/proxy"
)

func init() {
	logger.Debugf("tst")
}

func TestHttpCliGet(t *testing.T) {

	hc := NewHTTPClient("https://p.xgj.me:27035")
	//hc := NewHttpClient("ip://192.168.9.42")
	//hc := NewHttpClient("ip://[fe80::18c9:dfe5:85cd:e0c5]")
	//url := "https://play.google.com/store/apps/search?q=qq"
	url := "https://www.baidu.com"
	req, e := http.NewRequest(
		"Get",
		url,
		nil,
	)
	resp, e := hc.Do(req)
	if e != nil {
		log.Println(url)
		t.Fatal(e)
	}
	defer resp.Body.Close()
	buf, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		t.Fatal(e)
	}
	t.Log(string(buf))
}

func TestHttpCliPath(t *testing.T) {
	var key string
	var resURL string

	key = "encrypt key"
	resURL = "https://www.google.com"

	resURI, pErr := url.Parse(resURL)
	if pErr != nil {
		fmt.Print(pErr)
		return
	}

	fmt.Println("resUri:", resURI)

	path := resURI.EscapedPath()
	fmt.Println("path:", path)

	rawStr := fmt.Sprintf("%s%s%s", resURI, key, path)

	fmt.Println(rawStr)
}

func TestHttpCliGetWaiter(t *testing.T) {
	//hc := NewHttpClient("https://p.xgj.me:27035")
	var wg sync.WaitGroup
	hc := NewHTTPClient("ip://192.168.9.42")
	hc.Waiter = waiter.NewBurstLimitTick(time.Second, 3)
	time.Sleep(3 * time.Second)
	b := time.Now()

	for i := 0; i < 9; i++ {
		wg.Add(1)
		go func() {
			<-hc.Waiter.GetC()
			println("i:", i, time.Now().String())
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(time.Since(b))
}



func TestHttpCliGetNew(t *testing.T) {
	url := "https://www.google.com.hk/"
	proxys := []string{
		//"socks5://tnextday%40qq.com:10564757@s5.xgj.me:6124",
		"socks5://127.0.0.1:1086",
	}
	for _, proxyURL := range proxys {
		hc := NewHTTPClient(proxyURL)
		req, e := http.NewRequest(
			"GET",
			url,
			nil,
		)

		logger.Debugf("hc:%v", hc)
		resp, e := hc.Do(req)
		if e != nil {
			t.Fatal(e)
		}
		defer resp.Body.Close()
		buf, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			t.Fatal(e)
		}
		t.Log(string(buf))
	}

}

func TestSock5(t *testing.T) {
	//dialer, err := proxy.SOCKS5("tcp", "s5.xgj.me:6124",
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1086",
		//&proxy.Auth{
		//	User:"tnextday@qq.com",
		//	Password:"10564757",
		//},
		nil,
		&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		},
	)
	if err != nil {
		log.Fatalln("get dialer error", dialer)
	}
	httpTransport := &http.Transport{Dial: dialer.Dial}
	httpClient := &http.Client{Transport: httpTransport}
	resp, err := httpClient.Get("https://www.google.com.hk/")
	if err != nil {
		log.Fatalln(err)
	} else {
		defer resp.Body.Close()
		logger.Debugf("resp:%v", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%s\n", body)
	}

}
