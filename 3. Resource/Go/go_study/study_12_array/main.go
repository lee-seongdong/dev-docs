package main

import (
	"fmt"
)

func basic() {
	var arr [10]int // default value로 채워진다
	fmt.Println(arr)

	arr = [10]int{1} // 앞쪽부터 주어진 값으로 채워진다
	fmt.Println(arr)

	arr = [10]int{1: 10, 3: 30, 5: 50} // 인덱스:값 으로 채워진다
	fmt.Println(arr)

	arr2 := [...]int{10, 20, 30, 40, 50} // 배열 길이를 추론. 이 경우는 [5]int와 동일
	fmt.Println(arr2)

	const n = 10
	arr3 := [n]int{5} // 배열 선언 시 길이는 상수여야한다. 변수로 선언 시 에러
	fmt.Println(arr3)
}

func iterate() {
	nums := [...]int{10, 20, 30, 40, 50}
	nums[2] = 300

	for i := 0; i < len(nums); i++ {
		fmt.Println(nums[i])
	}

	t := [...]float64{24.0, 24.9, 27.8, 29.0}
	// range는 내장함수가 아니라 for 에서 사용되는 키워드이다
	for idx, value := range t {
		fmt.Println(idx, value)
	}
}

func copy() {
	a := [...]int{1, 2, 3, 4, 5}
	for i, v := range a {
		fmt.Printf("a[%d] = %d\n", i, v)
	}

	fmt.Println()
	b := [...]int{100, 200, 300, 400, 500}
	for i, v := range b {
		fmt.Printf("b[%d] = %d\n", i, v)
	}

	fmt.Println()
	b = a // 값을 복사하기 때문에 a값을 바꾸어도 영향받지 않음
	a[0] = 111
	for i, v := range b {
		fmt.Printf("b[%d] = %d\n", i, v)
	}
}

func multipleArr() {
	var a [2][5]int
	fmt.Println(a)

	a = [2][5]int{
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10}, // 닫는 중괄호가 다른줄에 있는 경우, 콤마는 필수
	}

	for _, a1 := range a {
		for _, a2 := range a1 {
			fmt.Print(a2, " ")
		}
		fmt.Println()
	}
}

func multipleArrCopy() {
	a := [2][5]int{
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10},
	}
	b := [2][5]int{
		{10, 20, 30, 40, 50},
		{60, 70, 80, 90, 100},
	}

	fmt.Println(a)
	fmt.Println(b)

	a = b
	b[0][0] = 111 // 값을 복사하기 때문에 b값을 바꾸어도 영향받지 않음

	fmt.Println(a)
	fmt.Println(b)
}

func main() {
	fmt.Println("#12 배열")
	// basic()
	// iterate()
	// copy()
	// multipleArr()
	multipleArrCopy()

	slice := []int{} // 배열 사이즈를 정하지 않으면 동적배열 구조체인 slice가 생성된다
	fmt.Println(slice)
}
