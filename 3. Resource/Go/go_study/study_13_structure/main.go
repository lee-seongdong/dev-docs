package main

import (
	"fmt"
	"unsafe"
)

type Student struct {
	Name  string
	Class int
	No    int
	Score float64
	pri   int
}

// 구조체를 필드로 사용사용할 수 있다
type School struct {
	students []Student
}

type TOPStudent struct {
	Student
	Name string
	top  int
}

type TestStruct1 struct {
	A int8 // 1byte
	B int  // 8byte
	C int8 // 1byte
	D int  // 8byte
	E int8 // 1byte
}

type TestStruct2 struct {
	A int8 // 1byte
	C int8 // 1byte
	E int8 // 1byte
	B int  // 8byte
	D int  // 8byte
}

func main() {
	fmt.Println("#13 구조체")
	// 여러 필드를 묶어서 사용하는 '타입'
	// 구조체의 역할 : 결합도는 낮추고 응집도는 높인다

	var s1 Student // 모든 필드값은 default value로 설정된다
	s1.Name = "홍길동"
	s1.Class = 1
	s1.No = 15
	s1.Score = 98.5
	s1.pri = 5
	fmt.Println(s1)

	var s2 Student = Student{"홍길동", 1, 15, 98.5, 3} // 필드값 초기화 (순서주의)
	fmt.Println(s2)

	var s3 Student = Student{Name: "홍길동", No: 3} // 특정 필드값만 초기화
	fmt.Println(s3)

	var school School = School{[]Student{s1, s2, s3}}
	fmt.Println(school)

	var topStudent TOPStudent // Student가 embedded Field로 설정되었다.
	topStudent.Student.Name = "홍길동"
	topStudent.Name = "아무개"
	topStudent.Student.No = 15
	fmt.Println(topStudent)

	fmt.Println("내장된 필드를 필드명 없이 접근 가능 ==>", topStudent.Student.No, topStudent.No)
	fmt.Println("필드명이 겹치는 embedded field는 타입명을 명시해야 접근 가능 ==>", topStudent.Student.Name, topStudent.Name)

	fmt.Println("#구조체 복사")
	fmt.Println(s1)
	fmt.Println(s2)
	s2 = s1
	s1.Name = "아무개" // 구조체 복사는 모든 필드값을 복사한다
	fmt.Println(s1)
	fmt.Println(s2)

	fmt.Println()
	fmt.Println("#구조체 메모리 & 메모리정렬 & 메모리 패딩")
	var i1 int
	var i2 int8
	var testStruct1 TestStruct1
	var testStruct2 TestStruct2
	fmt.Println("int size :", unsafe.Sizeof(i1), "byte, int8 size :", unsafe.Sizeof(i2), "byte")
	fmt.Println()
	fmt.Println("메모리는 8바이트 단위로 할당하기 때문에 8바이트 미만의 필드는 8바이트로 메모리를 패딩하여 할당한다")
	fmt.Println("struct size :", unsafe.Sizeof(testStruct1), "byte")
	fmt.Println()
	fmt.Println("8바이트 미만의 필드가 연속되면, 묶어서 할당한다.")
	fmt.Println("struct size :", unsafe.Sizeof(testStruct2), "byte")

}
