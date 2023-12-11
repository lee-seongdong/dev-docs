package main

import (
	"fmt"
)

func f1(num int) (int, bool) {
	return num, num > 5
}

func main() {
	// if 초기문; 조건문
	if n, f := f1(3); f {
		fmt.Println("f is true", n)
	} else {
		fmt.Println("f is false")
	}
}
