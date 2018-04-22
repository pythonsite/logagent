package main

import (
	"time"
	"sync/atomic"
	"github.com/astaxie/beego/logs"
)

type SecondLimit struct {
	unixSecond int64
	curCount int32
	limit int32
}

func NewSecondLimit(limit int32) *SecondLimit {
	secLimit := &SecondLimit{
		unixSecond:time.Now().Unix(),
		curCount:0,
		limit:limit,
	}
	return secLimit
}

func (s *SecondLimit) Add(count int) {
	sec := time.Now().Unix()
	if sec == s.unixSecond {
		atomic.AddInt32(&s.curCount,int32(count))
		return
	}
	atomic.StoreInt64(&s.unixSecond,sec)
	atomic.StoreInt32(&s.curCount, int32(count))
}

func (s *SecondLimit) Wait()bool {
	for {
		sec := time.Now().Unix()
		if (sec == atomic.LoadInt64(&s.unixSecond)) && s.curCount == s.limit {
			time.Sleep(time.Microsecond)
			logs.Debug("limit is running,limit:%d s.curCount:%d",s.limit,s.curCount)
			continue
		}

		if sec != atomic.LoadInt64(&s.unixSecond) {
			atomic.StoreInt64(&s.unixSecond,sec)
			atomic.StoreInt32(&s.curCount,0)
		}
		logs.Debug("limit is exited")
		return false
	}
}
