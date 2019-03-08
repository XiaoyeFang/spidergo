package spider

import (
	"fmt"
	"net/http"
	net"net/url"
)

func PostReq(url string) (code int) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	//http://192.168.0.18:10081
	purl, _ := net.Parse("http://192.168.0.18:10081")
	proxy := http.ProxyURL(purl)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: proxy,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		fmt.Println("2")
		fmt.Println("client.Do ", err)
		return
	}
	defer resp.Body.Close()

	if resp.Body == nil {
		fmt.Println("resp.Body == nil")
		return
	}
	fmt.Println(resp.Body)
	return resp.StatusCode
}
