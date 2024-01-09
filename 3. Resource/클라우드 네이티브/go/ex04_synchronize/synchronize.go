package synchronize

// '메모리를 공유하여 정보를 주고받지 말고, 정보를 주고받아 메모리를 공유하라'
// 대부분의 경우 채널을 통한 정보공유가 좋지만,
// 캐시나 대규모의 스테이트풀 구조체 접근에 대한 동기화에는 뮤텍스가 적합함

import (
	"sync"
)

type Resource1 struct {
	url        string
	polling    bool
	lastPolled int64
}

type Resources struct {
	data []*Resource1
	lock *sync.Mutex
}

func PollerWithLock(res *Resources) {
	for {
		res.lock.Lock()
		var r *Resource1

		for _, v := range res.data {
			if v.polling {
				continue
			}
			if r == nil || v.lastPolled < r.lastPolled {
				r = v
			}
		}

		if r != nil {
			r.polling = true
		}

		res.lock.Unlock()

		if r == nil {
			continue
		}

		res.lock.Lock()
		r.polling = false
	}
}

type Resource string

func PollerWithChan(in, out chan *Resource) {
	for r := range in {
		// logic
		out <- r
	}
}
