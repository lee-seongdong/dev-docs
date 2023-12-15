package main

import (
	"fmt"
)

type Data struct {
	value int
	data  [10]int
}

type User struct {
	name string
	age  int
}

func changeData(arg Data) {
	arg.value = 999
	arg.data[5] = 999
}

func changeData2(arg *Data) {
	(*arg).value = 999 // 편의상 arg.value도 동일하게 동작한다.
	(*arg).data[5] = 999
}

func testClearInstance() {
	u := &Data{}
	u.value = 30
	fmt.Println(u.value)
}

func newUser(name string, age int) *User {
	// 함수 내에서 사용된 변수는 함수가 끝나면서 사라진다.
	// C, C++에서는 &u 리턴 시 에러(dangling)

	// go 에서는 escape analysis를 통해 참조변수 유지
	// stack 영역에 생성된 지역변수는 사라지기 때문에, &u의 경우 탈출분석을 통해 Heap 공간에 생성된다
	var u = User{name, age}
	return &u
}

func main() {
	fmt.Println("#14 포인터")
	// 포인터 : 메모리 주소를 값으로 갖는 타입

	var a int
	var p1 *int
	p1 = &a  // &a : a가 저장된 주소값
	*p1 = 20 // *p : 주소 p에 저장된 값

	fmt.Println(p1, *p1, a)

	var p2 *int // 포인터는 default value : nil
	if p2 == nil {
		fmt.Println("nil", p2)
	}

	var data Data
	fmt.Println()
	fmt.Println("함수는 call by value 이므로, 인자의 값이 변경되지 않는다.")
	changeData(data)

	fmt.Println("data.value :", data.value)
	fmt.Println("data.data[5] :", data.data[5])

	fmt.Println()
	fmt.Println("함수에서 인자의 값을 변경하려면(call by reference), 인자의 주소를 넘겨야한다.")
	changeData2(&data)
	fmt.Println("data.value :", data.value)
	fmt.Println("data.data[5] :", data.data[5])

	fmt.Println()
	fmt.Println("#구조체 포인터 초기화")
	var data2 Data
	var p3 *Data = &data2 // 가리키는 값을 초기화 한 후 포인터 초기화
	fmt.Println("*p3", *p3)
	var p4 *Data = &Data{} // 가리키는 값과 함께 포인터 초기화
	fmt.Println("*p4", *p4)
	p5 := &Data{} // 타입추론으로 초기화
	fmt.Println("*p5", *p5)
	p6 := &data2
	fmt.Println("*p6", *p6)
	p7 := new(Data) // new 를 사용한 초기화. 기본값으로만 초기화 가능하다
	fmt.Println("*p7", *p7)

	// 인스턴스는 메모리에 생성된 데이터의 실체
	// 포인터를 사용해서 인스턴스를 가리킬 수 있다
	// 함수 호출 시 포인터 인수를 통해서 인스턴스를 전달받고, 인스턴스의 값을 변경할 수 있다.
	// 참조하는 포인터가 없는 인스턴스는 가비지컬렉터가 정리한다
	testClearInstance()

	//
	userPointer := newUser("aaa", 23)
	fmt.Println(*userPointer)
}
