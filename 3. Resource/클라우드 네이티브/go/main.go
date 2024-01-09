package main

import (
	"fmt"
	"lru_cache"
	"sync"
	"time"
)

func main() {
	// step1.Step1()
	// db_logger.Step2()
	lru_cache.Main()
	timer := time.NewTimer(2 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		for {
			select {
			case <-timer.C:
				fmt.Println("timer !")
			case <-ticker.C:
				fmt.Println("ticker !")
			}
		}
	}()
	wg.Wait()
}
