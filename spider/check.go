package spider

import (
	"errors"
	"fmt"
	"github.com/golang/glog"
	"net/http"
	"net/url"
	"time"
)

const (
	IDLECONNTIMEOUT = 5 * time.Second
)

func CheckAPKUpdate(name, apkurl, lastTime string) (err error) {
	glog.Infoln("check start")
	//设定Idetimeout和代理
	//http://172.16.8.8:30011
	purl, _ := url.Parse("http://172.16.8.8:31081")
	proxy := http.ProxyURL(purl)
	transport := new(http.Transport)
	transport.IdleConnTimeout = IDLECONNTIMEOUT
	transport.Proxy = proxy
	hc := http.Client{

		Transport: transport,
	}
	req, _ := http.NewRequest("GET", apkurl, nil)
	resp, err := hc.Do(req)
	if err != nil {
		glog.Errorf("url : %s,hc.Do err : %s \n", apkurl, err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		glog.Errorf("url : %s,resp.StatusCode : %d \n", apkurl, resp.StatusCode)
		err = errors.New(fmt.Sprintf("%d", resp.StatusCode))
		return
	}
	err = ParseAllVersion(resp.Body)
	if err != nil {
		return err
	}
	return nil
}
