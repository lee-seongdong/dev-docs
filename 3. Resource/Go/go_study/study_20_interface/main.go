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

type Reader interface {
	Read() (n int, err error)
	Close() error
}

type Writer interface {
	Writer() (n int, err error)
	Close() error
}

// 포함된 인터페이스. 내부에 선언된 인터페이스의 메서드 합집합을 가진 인터페이스
type ReadWriter interface {
	Reader
	Writer
}

type Attacker interface {
	Attack()
}

type Monster struct {
	Level int
}

func (m *Monster) Attack() {
	fmt.Println("monster Attack")
}

type Player struct {
	Level int
}

func (p *Player) Attack() {
	fmt.Println("player Attack")
}

func DoAttack(att Attacker) {
	if att != nil {
		fmt.Print("Do Attack! ")
		att.Attack()
		m, ok := att.(*Monster) // 인터페이스 타입 변환 결과와 성공 여부를 반환한다.
		if ok {
			fmt.Println(*m)
		} else {
			fmt.Println("this is not monster")
		}

		// 인터페이스 변환 및 사용의 관용적 표현
		// if 초기화문을 활용한다.
		if p, ok := att.(*Player); ok {
			fmt.Println(*p)
		}

	}
}

// 모든 타입을 받기위한 empty interface
func PrintVal(v interface{}) {
	// 인터페이스의 Concrete Type 변환
	switch t := v.(type) {
	case int:
		fmt.Printf("v is int %d\n", int(t))
	case float64:
		fmt.Printf("v is float %f\n", float64(t))
	case string:
		fmt.Printf("v is string %s\n", string(t))
	default:
		fmt.Printf("Not supported type %T:%v\n", t, t)
	}
}

func main() {
	fmt.Println("#20 인터페이스")

	student := Student{"tom", 20}
	// Student의 객체에 선언된 메서드 중 Stringer에 선언된 메서드 집합을 가져온다
	// 따라서, 메서드 집합을 가져오기 위해서는 Stringer에 선언한 인터페이스 모두를 Student에 메서드를 구현해야한다.
	var stringer Stringer = student
	fmt.Printf(stringer.String())

	fmt.Println("\n# 빈 인터페이스를 활용한 여러 타입 인자 전달")
	PrintVal(10)
	PrintVal(3.14)
	PrintVal("입력 테스트")
	PrintVal(student)

	p := Player{10}
	m := Monster{9}
	DoAttack(&p)
	DoAttack(&m)
}
