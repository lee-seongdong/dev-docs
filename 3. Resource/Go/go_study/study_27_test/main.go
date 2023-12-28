package main

import (
	"fmt"
)

var m map[int]int

func square(x int) int {
	return x * x
}

func fibonacci1(n int) int {
	if n < 0 {
		return 0
	}
	if n < 2 {
		return n
	}
	if m == nil {
		m = make(map[int]int)
	}
	if v, ok := m[n]; ok {
		return v
	}
	v := fibonacci1(n-1) + fibonacci1(n-2)
	m[n] = v
	return v
}

func fibonacci2(n int) int {
	if n < 0 {
		return 0
	}
	if n < 2 {
		return n
	}

	one := 1
	two := 0
	rst := 0
	for i := 2; i <= n; i++ {
		rst = one + two
		two = one
		one = rst
	}

	return rst
}

func main() {
	fmt.Println("#27 테스트와 벤치마크")
	fmt.Println("9 * 9 = ", square(9))
	fmt.Println("fibonacci1 : ", fibonacci1(13))
	fmt.Println("fibonacci1 : ", fibonacci2(13))
}
