package main

import (
	"fmt"
)

func add(a int, b int) int {
	return a + b
}

func divide(a, b int) (int, bool) {
	if b == 0 {
		return 0, false
	} else {
		return a / b, true
	}
}

func multiple(a, b int) (result int, success bool) {
	result = a * b
	success = true
	return
}

func fibo(n int) int {
	if n == 1 || n == 2 {
		return 1
	} else {
		return fibo(n-2) + fibo(n-1)
	}
}

func main() {
	// func 첫글자가 대문자인 경우 public으로 취급

	fmt.Println(add(3, 6))
	fmt.Println(divide(4, 2))

	a, success := divide(9, 3)
	fmt.Println(a, success)
	b, success := divide(10, 0) // 할당되지 않은 변수가 하나라도 있으면, 나머지 변수의 재 선언은 허용
	fmt.Println(b, success)
	fmt.Println(multiple(3, 3))
	fmt.Println(fibo(6))
}
