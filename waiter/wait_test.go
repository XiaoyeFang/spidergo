package waiter

import (
	"testing"
	//"fmt"
	"fmt"
	"sync"
	"time"
)

//func TestWaiter(t *testing.T)  {
//	w := NewRateLimit(1, 3)
//	//println(float64(time.Second) / 1)
//
//	//burstyRequests := make(chan string, 5)
//	//ipPoll := []string {"ip1", "ip2", "ip3",  "ip4", "ip5"}
//	//
//	//for _, v := range ipPoll{
//	//	burstyRequests <- v
//	//}
//	//
//	//defer  close(burstyRequests)
//	//for req := range burstyRequests {
//	//	w.GetThrottle(req)
//	//	fmt.Println("request", req, time.Now())
//	//}
//	for i:=0; i<5; i++ {
//		//println("ipgo:",time.Now().String())
//		go func(){
//			//println("ipaf:",time.Now().String())
//			w.GetThrottle("ip")
//			println("ipbf:",time.Now().String())
//		}()
//	}
//	for {
//		println("select", time.Now().String())
//		time.Sleep(time.Second)
//	}
//
//}
//
//func TestWaiter2(t *testing.T)  {
//	burstyLimiter := make(chan time.Time, 3)
//	for i := 0; i < 3; i++ {
//		burstyLimiter <- time.Now()
//	}
//
//	go func() {
//		for t := range time.Tick(time.Second * 2) {
//			burstyLimiter <- t
//		}
//	}()
//
//	ipPoll := []string {"ip1", "ip2", "ip3",  "ip4", "ip5"}
//
//	burstyRequests := make(chan string, 5)
//	for i := 0; i <= len(ipPoll)-1; i++ {
//		burstyRequests <- ipPoll[i]
//	}
//	close(burstyRequests)
//
//	for req := range burstyRequests {
//		<-burstyLimiter
//		fmt.Println("request", req, time.Now())
//	}
//
//}
//
func TestWaiter3(t *testing.T) {
	burstyLimiter := make(chan time.Time, 3)
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	go func() {
		for t := range time.Tick(time.Second * 2) {
			burstyLimiter <- t
		}
	}()

	ipPoll := []string{"ip1", "ip2", "ip3", "ip4", "ip5"}

	burstyRequests := make(chan string, 5)
	for i := 0; i <= len(ipPoll)-1; i++ {
		burstyRequests <- ipPoll[i]
	}
	close(burstyRequests)

	for req := range burstyRequests {
		<-burstyLimiter
		fmt.Println("request", req, time.Now())
	}

}

func TestWaiter4(t *testing.T) {
	var wg sync.WaitGroup
	w := NewBurstLimitTick(time.Second, 3)
	time.Sleep(3 * time.Second)
	b := time.Now()
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			<-w.GetC()
			println("t", time.Now().String())
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(time.Since(b))

}
