package main

import (
	"container/list"
	"container/ring"
	"fmt"
)

type Queue struct {
	v *list.List
}

type Stack struct {
	v *list.List
}

func (q *Queue) Push(val interface{}) {
	q.v.PushBack(val)
}

func (q *Queue) Pop() interface{} {
	front := q.v.Front()
	if front != nil {
		return q.v.Remove(front)
	} else {
		return nil
	}
}

func (s *Stack) Push(val interface{}) {
	s.v.PushBack(val)
}

func (s *Stack) Pop() interface{} {
	top := s.v.Back()
	if top != nil {
		return s.v.Remove(top)
	} else {
		return nil
	}
}

func NewQueue() *Queue {
	return &Queue{list.New()}
}

func NewStack() *Stack {
	return &Stack{list.New()}
}

func queueTest() {
	fmt.Println("# queue test")
	fmt.Println(`- list로 큐 구현
- 선입선출(FIFO)`)
	queue := NewQueue()
	fmt.Print("in : ")
	for i := 0; i < 10; i++ {
		fmt.Print(i, " ")
		queue.Push(i)
	}

	fmt.Print("\nout : ")
	v := queue.Pop()
	for v != nil {
		fmt.Print(v, " ")
		v = queue.Pop()
	}

	fmt.Println("")
	fmt.Println("")
}

func stackTest() {
	fmt.Println("# stack test")
	fmt.Println(`- list로 스택 구현
- 선입후출(FILO)`)
	stack := NewStack()
	fmt.Print("in : ")
	for i := 0; i < 10; i++ {
		fmt.Print(i, " ")
		stack.Push(i)
	}

	fmt.Print("\nout : ")
	v := stack.Pop()
	for v != nil {
		fmt.Print(v, " ")
		v = stack.Pop()
	}

	fmt.Println("")
	fmt.Println("")
}

func listTest() {
	fmt.Println("# list")
	fmt.Println(`- 요소 삽입, 삭제가 많은 경우 list가 유리
- 인덱스 접근이 많은 경우 array, slice가 유리
- 데이터 지역성으로 인해 일반적으로 요소수가 적은(1만개) 경우, 리스트보다 배열이 빠름
* 데이터 지역성 : 데이터가 인접해 있을 수록 캐시 성공률이 올라가고 성능이 증가하는 성질`)

	l := list.New()
	e4 := l.PushBack(4)
	e1 := l.PushFront(1)
	l.InsertBefore(3, e4)
	l.InsertAfter(2, e1)

	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Print(e.Value, " ")
	}

	fmt.Println()
	for e := l.Back(); e != nil; e = e.Prev() {
		fmt.Print(e.Value, " ")
	}
	fmt.Println()
	fmt.Println()
}

func ringTest() {
	fmt.Println("# ring test")
	fmt.Println(`- 일정 갯수만 사용하고, 오래된 요소가 지워져도 상관없는 경우 사용
- ex
  - 실행 취소 기능
  - 고정 크기 버퍼 기능
  - 리플레이 기능`)
	r := ring.New(5)
	n := r.Len()

	for i := 0; i < n; i++ {
		r.Value = 'A' + i
		r = r.Next()
	}

	fmt.Print("정방향 순회 : ")
	for i := 0; i < n; i++ {
		fmt.Printf("%c ", r.Value)
		r = r.Next()
	}

	fmt.Print("\n역방향 순회 : ")
	for i := 0; i < n; i++ {
		fmt.Printf("%c ", r.Value)
		r = r.Prev()
	}

	fmt.Println()
	fmt.Println()
}

type Product struct {
	Name  string
	Price int
}

func mapTest() {
	fmt.Println("# map test")
	fmt.Println(`- 종류
  - HashMap : key 순서 보장 되지 않음 (go 의 map은 HashMap이다)
  - SortedMap : key 순서 보장됩`)
	m := make(map[string]string)
	m["kim"] = "A"
	m["lee"] = "B"
	m["park"] = "C"

	delete(m, "lee")  // 삭제
	v, ok := m["lee"] // v는 default value로 설정되기 때문에, 존재여부는 반드시 ok 변수를 통해 확인해야한다.
	fmt.Println(v, ok)

	pm := make(map[int]Product)
	pm[16] = Product{"볼펜", 500}
	pm[46] = Product{"지우개", 2000}
	pm[78] = Product{"샤프심", 1000}
	pm[345] = Product{"샤프", 3000}
	pm[845] = Product{"샤프", 3000}

	for k, v := range pm {
		fmt.Println(k, v)
	}

	fmt.Println()
	fmt.Println()
}

const M = 10

func hash(d int) int {
	return d % M
}

func myHashMapTest() {
	fmt.Println("# 해쉬함수를 만들어 맵 구현")
	m := [M]string{}

	m[hash(23)] = "kim"
	m[hash(136)] = "Lee"
	m[hash(252)] = "Jang"
	m[hash(7)] = "Park"
	m[hash(13)] = "Tom" // 해쉬 충돌

	fmt.Printf("%d = %s\n", 23, m[hash(23)])
	fmt.Printf("%d = %s\n", 136, m[hash(136)])
	fmt.Printf("%d = %s\n", 252, m[hash(252)])
	fmt.Printf("%d = %s\n", 13, m[hash(13)])
	fmt.Printf("%d = %s\n", 11, m[hash(11)])

	fmt.Println()
	fmt.Println()
}

func main() {
	fmt.Println("#22 자료구조")
	listTest()
	queueTest()
	stackTest()
	ringTest()
	mapTest()
	myHashMapTest()
}
