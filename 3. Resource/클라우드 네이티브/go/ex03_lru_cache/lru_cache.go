package lru_cache

// lru는 글로벌 캐시로 유용하지만, 초당 700만 이상의 아주 높은 수준의 동시성이 필요한 경우 경합이 발생함

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru"
)

var cache *lru.Cache

func init() {
	cache, _ = lru.NewWithEvict(2, func(key, value interface{}) {
		fmt.Printf("Evicted: key = %v, value = %v\n", key, value)
	})
}

func Main() {
	cache.Add(1, "a")
	cache.Add(2, "b")

	fmt.Println(cache.Get(1))

	cache.Add(3, "c") // 캐시키 2가 제거됨
	fmt.Println(cache.Get(2))
}
