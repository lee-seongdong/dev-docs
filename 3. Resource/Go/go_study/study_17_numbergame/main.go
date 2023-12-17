package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var stdin = bufio.NewReader(os.Stdin)

func inputIntValue() (int, error) {
	var n int
	_, err := fmt.Scanln(&n)
	if err != nil {
		stdin.ReadString('\n')
	}

	return n, err
}

func main() {
	rand.Seed(time.Now().UnixNano())

	r := rand.Intn(100)
	count := 1
	for {
		fmt.Println("숫자를 입력하세요 :")
		n, err := inputIntValue()
		if err != nil {
			fmt.Println("숫자만 입력하세요")
		} else {
			if n > r {
				fmt.Println("DOWN")
			} else if n < r {
				fmt.Println("UP")
			} else {
				fmt.Println("정답. 시도한 횟수 : ", count)
				break
			}

			count++
		}
	}
}
