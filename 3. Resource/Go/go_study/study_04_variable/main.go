package main

import "fmt"

func variable1() {
	// 변수 선언하는 방법 4가지
	var a int = 10 // 초기값 설정
	var c int      // default value (zero value로 설정됨)
	var b = 10     // 타입 추론을 위해 대입값 필수
	msg := "hello" // 타입 추론을 위해 대입값 필수

	/*
		타입별 기본값
		- 정수 : 0
		- 실수 : 0.0
		- boolean : false
		- string : ""
		- 나머지 : nil
	*/

	fmt.Println(a, b, c, msg)
}

func variable2() {
	// 연산의 각 항목은 타입이 반드시 같아야 한다.
	// 타입은 물론이고 크기도 같아야한다 (alias로 부여한 별칭은 다른 타입으로 간주한다)
	var a int32 = 1
	var b int8 = 2
	var c float32 = 2.2
	var d int = int(c)
	fmt.Println(a+int32(b), c, d)

	type myInt int8 // 타입 별칭
	var e myInt = 5
	var f int8 = int8(e)
	fmt.Println(f)
}

func variable3() {
	// overflow test
	var a int16 = 3456 // 2byte 정수
	var b int8 = int8(a)
	fmt.Println(a, b)
}

func variable4() {
	fmt.Println("float 표현 범위 테스트")
	// float32 정밀도 : 7자리
	// float64 정밀도 : 15자리
	var a float32 = 1234.523
	var b float32 = 3456.123
	var c float32 = a * b
	var d = c * 3
	fmt.Println(a, b, c, d)
}

func main() {
	//variable1()
	//variable2()
	//variable3()
	variable4()
}
