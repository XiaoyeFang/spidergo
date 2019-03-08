package spider

import (
	"testing"
)

func TestPostReq(t *testing.T) {

//	sem := make(chan int, 100)
//	defer close(sem)
//	for i := 0; i < 5000000; i++ {
//		jsonStr := []byte(`{
//"fileId": "b/apk/Y29tLmFuaXBsZXguZmF0ZWdyYW5kb3JkZXJfMTM3Xzk0NjQ0MTY` + fmt.Sprintf("%d", i) + `",
//"expireSeconds": 172800,
//"clientCountry": "SY",
//"fileName": "Fate Grand Order` + fmt.Sprintf("%d", i) + `",
//"hotFile": false}`)
//
//		url := "http://192.168.9.111:9090/v1/fs_router/get_url"
//
//		sem <- 1
//		go func(i int) {
//		t.Log(i)
//		code := PostReq(url, jsonStr)
//		if code == 200 {
//			t.Log("success")
//		}
//			<-sem
//		}(i)
//
//	}
code :=PostReq("https://graph.facebook.com/1997670463844922/picture?type=normal")
	//if err != nil {
		t.Errorf(" %d \n", code)
	//}

}
