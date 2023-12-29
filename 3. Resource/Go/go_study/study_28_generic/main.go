package main

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

type Stringer interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
	String() string
}

// [T any] : 타입 파라미터
func print1[T any](a T) {
	fmt.Println(a)
}

// 타입 파라미터에 인터페이스도 지정이 가능
// 파라미터의 타입이 인터페이스인 경우, 인터페이스의 함수만 사용이 가능함
// 타입파라미터로 인터페이스를 받으면, 실제 전달받은 타입으로 동작함
func Print2[T Stringer](a T) {
	fmt.Println(a)
}

func PrintMin[T Stringer](a, b T) {
	if a < b {
		fmt.Println(a.String())
	} else {
		fmt.Println(b.String())
	}
}

func min1(a, b interface{}) interface{} {
	// 모든 타입이 비교연산을 지원하지 않기 때문에 에러 발생함
	// if a < b {
	// 	return a
	// }
	return b
}

// 자주 사용하는 제네릭은 제공되고있다
// 제네릭 타입 선언
type Inteager interface {
	~int | ~int16 | ~int32 | uint8 | uint16 | uint32
	// ~ 를 사용하면, 해당 타입에 별칭을 부여한 타입에도 적용한다는 의미
}

type MyInt int

func (m MyInt) String() string {
	return strconv.Itoa(int(m))
}

type Float interface {
	float32 | float64
}

type Numeric interface {
	Inteager | Float
}

// 제네릭 타입
type Node[T any] struct {
	val  T
	next *Node[T]
}

// 타입 파라미터에 타입 제한을 두면 빈 인터페이스의 문제를 해결할 수 있음 (가능한 연산을 미리 알 수 있다)
// 언박싱 없이 값을 바로 사용 가능하다
func min2[T int | int16 | int32 | uint8 | uint16 | uint32 | float32 | float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func min3[T Numeric](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func min4[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Map[F, T any](s []F, f func(F) T) []T {
	rst := make([]T, len(s))
	for i, v := range s {
		rst[i] = f(v)
	}
	return rst
}

func NewNode[T any](v T) *Node[T] {
	return &Node[T]{
		val: v,
	}
}

func (n *Node[T]) Push(v T) *Node[T] {
	node := NewNode(v) // 타입 추론으로 [T] 생략 가능
	n.next = node
	return node
}

func main() {
	fmt.Println("#28 제네릭")
	fmt.Println("제네릭은 Go 1.18에 도입된 개념")
	print1("test")
	print1(1)

	min3('a', 'b')
	min3(1.1, 2.2)
	min3(1, 3)

	min4('a', 'b')
	min4(1.1, 2.2)
	min4(1, 3)

	var m1 MyInt = 10
	var m2 MyInt = 20
	PrintMin(m1, m2)

	doubled := Map([]int{1, 2, 3}, func(v int) int {
		return v * 2
	})

	fmt.Println(doubled)

	uppered := Map([]string{"hello", "world", "go lang"}, func(v string) string {
		return strings.ToUpper(v)
	})
	fmt.Println(uppered)

	n1 := NewNode[int](1) // 타입추론으로 [int] 생략 가능
	n2 := NewNode(2)
	n3 := NewNode(3)
	n1.next = n2
	n2.next = n3

	n3.Push(4).Push(5).Push(6)

	next := n1
	for next != nil {
		fmt.Print(next.val, " ")
		next = next.next
	}
}
