package main

import (
	"fmt"
)

// 인터페이스도 타입이다
type DuckInterface interface {
	Fly()
	// Fly(distance int) int // 메서드명은 겹치면 안된다. 따라서 메소드 오버로드가 불가능하다.
	Walk(distance int) int
}

type Stringer interface {
	String() string
	//String3() string
}

type Student struct {
	name string
	age  int
}

func (s Student) String() string {
	return fmt.Sprintf("Hello! name :%s, age : %d", s.name, s.age)
}

func (s Student) String2() string {
	return fmt.Sprintf("Hello2! name :%s, age : %d", s.name, s.age)
}

func main() {
	fmt.Println("#20 인터페이스")

	student := Student{"tom", 20}
	// Student의 객체에 선언된 메서드 중 Stringer에 선언된 메서드 집합을 가져온다
	// 따라서, 메서드 집합을 가져오기 위해서는 Stringer에 선언한 인터페이스 모두를 Student에 메서드를 구현해야한다.
	var stringer Stringer = student
	fmt.Printf(stringer.String())
}
