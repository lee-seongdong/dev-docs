package main

import (
	"fmt"
	"os"
)

// ...type : 가변길이 인자타입
func sum1(nums ...int) int {
	sum := 0
	fmt.Printf("nums 타입: %T\n", nums) // 가변길이 인자타입은 슬라이스로 동작한다
	for _, v := range nums {
		sum += v
	}
	return sum
}

func sum2(nums []int) int {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	return sum
}

// 함수 실행 종료 직전에 실행을 보장하는 키워드 (finally와 유사함)
// 주로 os 자원 반납에 사용함
// defer는 스택으로 기록되므로, 마지막 호출한 defer부터 실행됨
func fileTest() {
	f, err := os.Create("test.txt")
	if err != nil {
		fmt.Println("Failed to create a file", err)
		return
	}

	defer fmt.Println("반드시 호출됨")
	defer f.Close()
	defer fmt.Println("파일을 닫습니다.")

	fmt.Println("파일에 Hello World를 씁니다.")
	fmt.Fprintln(f, "Hello World")
}

// 함수 타입은 함수 시그니쳐로 표현
// add 함수의 타입 : func (int, int) int
func add(a, b int) int {
	return a + b
}

func mul(a, b int) int {
	return a * b
}

type opFunc func(int, int) int // 함수타입도 별칭으로 사용 가능
func getOperator(op string) func(int, int) int {
	if op == "+" {
		return add
	} else if op == "*" {
		return mul
	} else {
		return nil
	}
}

func captureLoop1() {
	f := make([]func(), 3)
	fmt.Println("ValueLoop1")
	for i := 0; i < 3; i++ {
		// i 포인터는 하나로 재사용 되기 때문에, 최종적으로 fmt.Println(3) 으로 실행된다.
		f[i] = func() { fmt.Println(i) }
	}

	for i := 0; i < 3; i++ {
		f[i]()
	}
}

func captureLoop2() {
	f := make([]func(), 3)
	fmt.Println("ValueLoop2")
	for i := 0; i < 3; i++ {
		v := i // v 포인터는 매번 재 생성 되기 때문에, 0, 1, 2가 출력된다.
		f[i] = func() { fmt.Println(v) }
	}

	for i := 0; i < 3; i++ {
		f[i]()
	}
}

type Writer func(string)
type WriterInterface interface {
	Write(string)
}

func writeHello(writer Writer) {
	writer("Hello World")
}

func writeHello2(writer WriterInterface) {
	writer.Write("Hello World")
}

func dependencyInjection() {
	f, err := os.Create("test.txt")
	if err != nil {
		fmt.Println("Failed to create a file")
		return
	}

	defer f.Close()

	writeHello(func(msg string) {
		fmt.Fprintln(f, msg)
	})
	writeHello(func(msg string) {
		fmt.Println(msg)
	})
}

func main() {
	fmt.Println("#21 함수고급")
	fmt.Println("# 가변 길이 인자 타입")
	fmt.Println(sum1(1, 2, 3, 4, 5))
	fmt.Println(sum2([]int{10, 12, 33, 44, 55}))

	fmt.Println("\n# defer 지연실행")
	fileTest()

	fmt.Println("\n# 함수타입 변수")
	var addOperator opFunc = getOperator("+")
	fmt.Println(addOperator(1, 3))
	fmt.Println(getOperator("*")(1, 3))

	fmt.Println("\n# 함수 리터럴")

	literalF1 := func(a, b int) int {
		return a - b
	}
	fmt.Println(literalF1(8, 1))

	fmt.Println("일반 함수는 상태를 가질 수 없지만, 함수 리터럴은 내부 상태를 가질 수 있다. (Capture)")
	i := 0
	literalF2 := func() {
		i++ // i 의 레퍼런스를 capture 해서 가진다
	}
	literalF2()
	literalF2()
	literalF2()
	fmt.Println(i)
	captureLoop1()
	captureLoop2()

	fmt.Println("\n# 의존성 주입")
	dependencyInjection()
}
