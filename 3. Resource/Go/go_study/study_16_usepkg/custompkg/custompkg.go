package custompkg

import (
	"fmt"
	"go_study/study_16_usepkg/exinit"
)

type Student struct {
	Name  string // public
	Age   int    // public
	score int    // private!
}

// 네이밍의 첫글자가 대문자인 경우, 패키지 내의 타입, 전역변수, 상수, 함수 등을 외부 공개해서 사용 가능
func PrintCustom() {
	fmt.Println("This is custom package!")
	exinit.PrintD()
}

// private function
func printCustom2() {
	fmt.Println("this is private package")
}
