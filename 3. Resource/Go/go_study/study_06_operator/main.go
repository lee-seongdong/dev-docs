package main

import (
	"fmt"
	"math"
)

// 실수 연산 오차를 피하는 방법 중 하나
// 실수 연산 시 오차를 방지하려면, 정수로 바꿔서 연산을 하자
func equal(a, b float64) bool {
	return math.Nextafter(a, b) == b
}

func main() {
	var a float64 = 0.1
	var b float64 = 0.2
	var c float64 = 0.3

	fmt.Printf("%f + %f == %f, %v", a, b, c, a+b == c)
	fmt.Println(a + b)
	fmt.Printf("%0.18f == %0.18f : %v\n", c, a+b, equal(a+b, c))
}
