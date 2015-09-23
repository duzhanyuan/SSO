package monitor

import (
	"container/list"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxLen = 10000
)

var (
	count   int64
	cntList *list.List
	mu      sync.Mutex
)

func Work() {
	cntList = list.New()
	for {
		time.Sleep(time.Second * 5)
		mu.Lock()
		cntList.PushBack(count)
		if len(cntList) > maxLen {
			it := cntList.Front()
			cntList.Remove(it)
		}
		count = 0
		mu.Unlock()
	}
}

func IncrCount() {
	atomic.AddInt64(&count, 1)
	fmt.Println("now count is: ", count)
}
