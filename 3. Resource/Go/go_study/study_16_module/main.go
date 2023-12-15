package main

import (
	"fmt" // 코드에서 사용하는 패키지 명은 최 하단 경로명(rand)
	"html/template"
	"math/rand"
	_ "strings"                  // 사용하지 않는 패키지를 import 하는 경우 오류발생. 패키지 초기화(init 함수)에 따른 부가효과를 위해서는 _ 로 별칭을 주어 오류를 방지한다
	textTemplate "text/template" // 겹치는 패키지명은 별칭으로 구분
)

func main() {
	fmt.Println("#16 모듈과 패키지")
	// 모듈 : 1개 이상의 패키지의 모음
	// 패키지 : 코드를 묶는 단위
	// 프로그램(실행파일) : 실행 시작지점(main 함수)을 포함한 패키지
	fmt.Println(rand.Int())
	template.New("foo").Parse(`{{define "T"}}Hello`)
	textTemplate.New("foo").Parse(`{{define "T"}}Hello`)
}
