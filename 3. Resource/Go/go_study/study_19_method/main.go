package main

import (
	"fmt"
	"time"
)

type Student struct {
	age  int
	name string
}

type User struct {
	Name string
	Age  int
}

type VIPUser struct {
	User     // embedded field
	VIPLevel int
}

func (u User) string() string {
	return fmt.Sprintf("%s, %d", u.Name, u.Age)
}

// 메서드 오버라이딩이 아니다
func (v VIPUser) string() string {
	return fmt.Sprintf("%s, %d", v.Name, v.VIPLevel)
}

func (v VIPUser) vipLevel() int {
	return v.VIPLevel
}

// 함수
func exFunc(s *Student) int {
	return s.age
}

// 메서드. 리시버 타입에 속한 함수
// s Student : 리시버. 리시버는 해당 지역의 타입에 대해서만 설정 가능
// 객체 = 데이터 + 기능
// object = state + function
func (s Student) exMethod() int {
	return s.age
}

type account struct {
	balance   int
	firstName string
	lastName  string
}

func (a1 *account) widthrawPointer(amount int) {
	a1.balance -= amount
}

func (a2 account) widthrawValue(amount int) {
	// a2는 호출하는 account type과 별도의 인스턴스이다.
	a2.balance -= amount
}

func (a2 account) widthrawValue2(amount int) account {
	// a2는 호출하는 account type과 별도의 인스턴스이다.
	a2.balance -= amount
	return a2
}

// 포인터 타입 메서드와 값 타입 메서드의 차이
// 포인터타입 메서드: 필드가 변경되어도, 동일한 객체로 볼 수 있는 경우 사용
// 값타입 메서드   : 필드가 변경되면, 동일한 객체로 볼 수 없는 경우 사용

func main() {
	fmt.Println("#19 메서드")

	student1 := Student{18, "홍길동"}
	fmt.Println(exFunc(&student1))
	fmt.Println(student1.exMethod())

	var mainA *account = &account{100, "joe", "kim"}
	mainA.widthrawPointer(30)
	fmt.Println(mainA.balance)

	mainA.widthrawValue(30)    // (*mainA).widthrawValue 와 동일
	fmt.Println(mainA.balance) // mainA는 변경 없음

	*mainA = mainA.widthrawValue2(30)
	fmt.Println(mainA.balance)

	// go 에는 생성자, 소멸자가 없다.
	// 생성자가 없기 때문에, 객체를 생성하는 방법을 강제할 수 없다.
	// 객체를 생성하는 메서드를 사용하도록 유도해야한다.
	var t1 = time.Timer{}
	var t2 = time.NewTimer(time.Second)
	fmt.Println(t1)
	fmt.Println(t2)

	vip := VIPUser{User{"홍길동", 10}, 5}
	fmt.Println(vip.string())
	fmt.Println(vip.User.string())
}
