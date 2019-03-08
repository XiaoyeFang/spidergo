package ippoll

import (
	"testing"
	"time"
	"sync"
	"io/ioutil"
	"net/http"
	"fmt"
)

func TestGetIp(t *testing.T)  {
	ipList := getIP()
	t.Logf("iplist:%v", ipList)

}

func TestHttpCliIpPoll(t *testing.T) {
	fmt.Printf("iplist:%v\n", IPPoll)
	t.Logf("ippoll:%v", IPPoll)

	var wg sync.WaitGroup
	b := time.Now()
	for index := range IPPoll {
		t.Logf("index:%v  hc:%v", index, IPPoll[index])
		wg.Add(1)
		go func(index int) {
			url := "https://www.baidu.com"
			req, e := http.NewRequest(
				"Get",
				url,
				nil,
			)
			if e != nil {
				t.Error(e)
			}
			resp, e := IPPoll[index].Do(req)
			if e != nil {
				t.Fatal(e)
			}

			defer resp.Body.Close()
			buf, e := ioutil.ReadAll(resp.Body)
			if e != nil {
				t.Fatal(e)
			}
			if len(buf) != 0 {
				t.Log("client:", IPPoll[index].ProxyAddr, " ", "ok:", len(buf), string(buf[:30]))
				println()
				println()
			}
			wg.Done()

		}(index)
	}

	wg.Wait()
	fmt.Println("time:", time.Since(b))
}

func TestGetHc(t *testing.T) {
	hc := GetHc()
	index := GetIPIndex()
	t.Log(index, hc)
	hc = GetHc()
	index = GetIPIndex()
	t.Log(index, hc)
}