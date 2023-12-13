package main

import (
	"fmt"
	"time"
)

func basicLoop() {
	// 초기문; 조건문; 후처리
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
}

func infinityLoop() {
	i := 0
	for {
		if i == 10 {
			break
		}

		time.Sleep(time.Second)
		fmt.Println(i)
		i++
	}
}

func nestedLoop() {
OuterFor: // 가능하면 쓰지 않는 것이 좋다
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			for k := 0; k < 10; k++ {
				fmt.Println(i, j, k)
				if k == 3 {
					break OuterFor
				}
			}
		}
	}

	fmt.Println("break outer")
}

func main() {
	basicLoop()
	// infinityLoop()
	nestedLoop()
}
