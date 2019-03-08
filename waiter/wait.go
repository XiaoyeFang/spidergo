package waiter

import (
	"sync"
	"time"
)

//BurstLimitTick contrl fz
type BurstLimitTick struct {
	tick        *time.Ticker
	burstLimit  int
	burstLimitC chan time.Time
}

//RateLimit to limit rate
type RateLimit struct {
	throttles  map[string]*BurstLimitTick
	rate       time.Duration
	burstLimit int
	lock       sync.RWMutex
}

//NewRateLimit new
func NewRateLimit(rqs float64, burstLimit int) *RateLimit {
	return &RateLimit{
		throttles:  make(map[string]*BurstLimitTick),
		rate:       time.Duration(float64(time.Second) / rqs),
		burstLimit: burstLimit,
	}
}

func (rrm *RateLimit) RequestWait(key string) {
	<-rrm.GetThrottle(key)
}

func (rrm *RateLimit) GetThrottle(key string) <-chan time.Time {
	rrm.lock.RLock()
	t, ok := rrm.throttles[key]
	rrm.lock.RUnlock()
	if ok {
		return t.GetC()
	}
	rrm.lock.Lock()
	defer rrm.lock.Unlock()
	t = NewBurstLimitTick(rrm.rate, rrm.burstLimit)
	rrm.throttles[key] = t
	return t.GetC()
}

func NewBurstLimitTick(rate time.Duration, burstLimit int) *BurstLimitTick {
	blt := &BurstLimitTick{
		tick:       time.NewTicker(rate),
		burstLimit: burstLimit,
	}
	if burstLimit > 0 {
		blt.burstLimitC = make(chan time.Time, burstLimit)
		go func() {
			for t := range blt.tick.C {
				select {
				case blt.burstLimitC <- t:
				default:
				}
			} // exits after blt.tick.Stop()
		}()
	}
	return blt
}

func (blt *BurstLimitTick) Stop() {
	blt.tick.Stop()
}

func (blt *BurstLimitTick) GetC() <-chan time.Time {
	if blt.burstLimit > 0 {
		return blt.burstLimitC
	} else {
		return blt.tick.C
	}
}
