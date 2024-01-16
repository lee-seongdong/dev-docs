package main

import (
	"context"
	"hexarch/core"
	"hexarch/frontend"
	"hexarch/transact"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func A() {
	tl, err := transact.NewTransactionLogger("postgres")
	if err != nil {
		log.Fatal(err)
	}

	store := core.NewKeyValueStore[string](tl)
	store.Restore()

	fe, err := frontend.NewFrontEnd("rest")
	if err != nil {
		log.Fatal(err)
	}

	err = fe.Start(store)
	if err != nil {
		log.Fatal(err)
	}
}

// 스로틀링 하고자 하는 대상 함수
type Effector func(context.Context) (string, error)

// Throttled는 Effector의 래퍼
type Throttled func(context.Context, string) (bool, string, error)

type bucket struct {
	tokens uint
	time   time.Time
}

// Throttle은 Effector를 매개변수로 받고,
// Throttled 함수를 UID별로 만들어진 토큰 버킷의 최대 용량과 함께 반환하며,
// 최대 용량은 매 d마다 리필토큰 비율에 따라 다시 채워짐
func Throttle(e Effector, max uint, refill uint, d time.Duration) Throttled {
	buckets := map[string]*bucket{}

	return func(ctx context.Context, uid string) (bool, string, error) {
		b := buckets[uid]
		if b == nil {
			buckets[uid] = &bucket{tokens: max - 1, time: time.Now()}
			str, err := e(ctx)
			return true, str, err
		}

		// 이전 요청 후 흐른 시간을 기준으로 남은 토큰 계산
		refillInterval := uint(time.Since(b.time) / d)
		tokensAdded := refill * refillInterval
		currentTokens := b.tokens + tokensAdded

		// 토큰이 부족한 경우 false 리턴
		if currentTokens < 1 {
			return false, "", nil
		}

		// 버킷이 리필되면 시간 다시 설정
		if currentTokens > max {
			b.time = time.Now()
			b.tokens = max - 1
		} else { // 아닌 경우, 가장 최근에 토큰이 추가된 시점 확인
			deltaTokens := currentTokens - b.tokens
			deltaRefills := deltaTokens / refill
			deltaTime := time.Duration(deltaRefills) * d

			b.time = b.time.Add(deltaTime)
			b.tokens = currentTokens - 1
		}

		str, err := e(ctx)

		return true, str, err
	}
}

var throttle = Throttle(getHostName, 1, 1, time.Second)

func getHostName(ctx context.Context) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	return os.Hostname()
}

func throttleHandler(w http.ResponseWriter, r *http.Request) {
	ok, hostname, err := throttle(r.Context(), r.RemoteAddr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Too many requested", http.StatusTooManyRequests)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hostname))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hostname", throttleHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
