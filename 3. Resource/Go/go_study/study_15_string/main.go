package main

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func toUpper1(str string) string {
	var result string
	for _, c := range str {
		if 'a' <= c && c <= 'z' {
			result += string('A' + c - 'a')
		} else {
			result += string(c)
		}
	}

	return result
}

func toUpper2(str string) string {
	var builder strings.Builder
	for _, c := range str {
		if 'a' <= c && c <= 'z' {
			builder.WriteRune('A' + c - 'a')
		} else {
			builder.WriteRune(c)
		}
	}

	return builder.String()
}

func main() {
	fmt.Println("#15 문자열")
	poet1 := "죽는 날까지 하늘을 우러러\n한 점 부끄럼이 없기를"
	poet2 := `죽는 날까지 하능을 우러러
한 점 부끄럼이 없기를` // 탭, 특수문자까지도 문자로 인식함

	fmt.Println(poet1)
	fmt.Println(poet2)

	// golang은 기본적으로 utf-8로 문자 표현
	// 문자열은 바이트 단위로 배열을 구성함
	str := "Hello 월드"
	fmt.Println("# 문자열 인덱스로 순회")
	for i := 0; i < len(str); i++ {
		fmt.Printf("타입: %T, 값: %d, 문자: %c\n", str[i], str[i], str[i])
	}

	fmt.Println("\n# 문자열 range로 순회")
	for _, v := range str {
		fmt.Printf("타입: %T, 값: %d, 문자: %c\n", v, v, v)
	}

	arr := []rune(str) // rune : int32의 별칭
	fmt.Println("\n# 문자열 rune으로 변환 후 순회")
	for i := 0; i < len(arr); i++ {
		fmt.Printf("타입: %T, 값: %d, 문자: %c\n", arr[i], arr[i], arr[i])
	}

	str1 := "Hello"
	str2 := "hello"
	str3 := "Hello"
	str4 := "World"
	fmt.Println("\n# 문자열 비교")
	fmt.Println("str1 == str2", str1 == str2)
	fmt.Println("str1 == str3", str1 == str3)
	fmt.Println("str2 == str3", str2 == str3)
	fmt.Println("str1 < str4", str1 < str4)

	str5 := ""
	str5 = str1
	// 문자열의 길이가 다른데(메모리 공간이 다른데) 대입이 가능한 이유
	// 값의 시작 포인터와 바이트 길이를 저장하는 구조이기 때문
	// type StringHeader struct {
	// 	Data uintptr
	// 	Len int
	// }

	fmt.Println("\n# 문자열의 구조와 대입연산")
	stringPointer1 := unsafe.Pointer(&str1)
	stringPointer5 := unsafe.Pointer(&str5)

	stringHeader1 := (*reflect.StringHeader)(stringPointer1)
	stringHeader2 := (*reflect.StringHeader)(stringPointer5)

	fmt.Println(stringHeader1)
	fmt.Println(stringHeader2)

	fmt.Println("\n# 문자열은 불변이다")
	// 문자열이 불변인 이슈
	// 문자열이 가변인 경우, 값이 변경되면 이를 가리키는 참조변수들도 모두 영향을 받는다
	// 따라서 안정성을 위해서는 대입할 때 마다 새로운 메모리 공간을 할당해야한다.

	// 문자열이 불변인 경우, 하나의 참조변수에서 값이 변경되어도(재할당) 나머지 참조변수에서는 영향을 받지 않기 때문에 안정적이다.
	// 이 경우 문자열을 대입할 때는 공유, 변경될때 복사하는 방식으로 동작한다.
	// 내생각 : 문자열의 변경보다 할당의 비율이 높아서 문자열 변경 시 메모리 할당하는 방식을 채택한 것이 아닐까? (효율적)
	// str5[1] = 'a' // 문자열의 일부를 바꾸는것은 불가능하다

	var str6 string = "Hello World"
	var slice []byte = []byte(str6)

	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&str6))
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&slice))

	fmt.Printf("str: \t%x\n", stringHeader.Data)
	fmt.Printf("slice: \t%x\n", sliceHeader.Data)

	fmt.Println("\n# 문자열의 연산은 새로운 문자열을 생성한다")
	fmt.Println(toUpper1("Hello World"))
	fmt.Println(toUpper2("Hello World"))
}
