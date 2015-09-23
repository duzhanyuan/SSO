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

func GetAllData() []int64 {
	data := []int64{}
	mu.Lock()
	for it := cntList.Front(); it != nil; it = it.Next() {
		data = append(data, it.Value)
	}
	mu.Unlock()
	return data
}

/*func GetData(seconds int64) []int64 {*/
//data := []int64{}
//seconds = seconds / 5
//for it := cntList.Back(); it != nil; it = it.Prev() {
//data = append(data, it.Value)
//seconds = seconds - 1
//if seconds == 0 {
//break
//}
//}
//return data
/*}*/
