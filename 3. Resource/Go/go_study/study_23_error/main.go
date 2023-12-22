package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type PasswordError struct {
	Len        int
	RequireLen int
}

// custom error 객체로 사용하기 위한 인터페이스
func (err PasswordError) Error() string {
	return "암호길이가 짧습니다."
}

func RegisterAccount(name, password string) error {
	if len(password) < 8 {
		return PasswordError{len(password), 8}
	}

	return nil
}

func Sqrt(f float64) (float64, error) {
	if f < 0 {
		// return 0, fmt.Errorf("제곱근은 양수여야합니다. f: %g", f)
		return 0, errors.New("제곱근은 양수여야 합니다.")
	}
	return math.Sqrt(f), nil
}

func ReadFile(filename string) (string, error) {
	file, error := os.Open(filename)
	if error != nil {
		return "", error
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	line, _ := reader.ReadString('n')
	return line, nil
}

func WriteFile(filename string, line string) error {
	file, error := os.Create(filename)
	if error != nil {
		return error
	}

	defer file.Close()

	_, error = fmt.Fprintln(file, line)
	return error
}

func readNextInt(scanner *bufio.Scanner) (int, int, error) {
	if !scanner.Scan() {
		return 0, 0, fmt.Errorf("Failed to scan")
	}

	word := scanner.Text()
	number, err := strconv.Atoi(word) // "24" -> 24
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to convert number")
	}

	return number, len(word), nil
}

func readEq(eq string) {
	rst, err := MultipleFromString(eq)
	if err != nil {
		fmt.Println(err)
		var numError *strconv.NumError
		if errors.As(err, &numError) {
			fmt.Println("numberError", numError)
		}
	} else {
		fmt.Println(rst)
	}
}

func MultipleFromString(str string) (int, error) {
	scanner := bufio.NewScanner(strings.NewReader(str))
	scanner.Split(bufio.ScanWords)

	pos := 0
	a, n, err := readNextInt(scanner)
	if err != nil {
		// %w 포매터로 에러 래핑 가능
		// errors.As() 함수로 래핑된 에러객체 접근 가능
		return 0, fmt.Errorf("Failed to readNextInt(), pos:%d, err:%w", pos, err)
	}

	pos += n + 1
	b, n, err := readNextInt(scanner)
	if err != nil {
		return 0, fmt.Errorf("Failed to readNextInt(), pos:%d, err:%w", pos, err)
	}

	return a * b, nil
}

func divide(a, b int) {
	if b == 0 {
		panic("b는 0 일 수 없습니다.")
	}
	fmt.Printf("%d / %d = %d\n", a, b, a/b)
}

func g() {
	divide(9, 3)
	divide(9, 0)
}

func f() {
	fmt.Println("f() 시작")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic 복구 - ", r)
		}
	}()

	g()
	fmt.Println("f() 끝")
}

func fileReadWriteTest() {
	fmt.Println("# 파일 입출력 에러 헨들링 테스트")
	const filename = "data.txt"
	line, err := ReadFile(filename)
	if err != nil {
		err = WriteFile(filename, "this is writeFile")
		if err != nil {
			fmt.Println("파일 생성 실패")
			return
		}

		line, err = ReadFile(filename)
		if err != nil {
			fmt.Println("파일 읽기 실패")
			return
		}
	}

	fmt.Println("파일내용 :", line)
	fmt.Println()
	fmt.Println()
}

// 사용자 에러 생성 방법
// 1. fmt.Errorf(formatter string, ...interface{}) error
// 2. errors.New(text string) error
func customErrorTest() {
	fmt.Println("# 커스텀 에러 생성 테스트")
	sqrt, err := Sqrt(-2)
	if err != nil {
		fmt.Println("value, Error :", sqrt, err)
	}

	fmt.Println()
	fmt.Println()
}

func customErrorTest2() {
	fmt.Println("# 커스텀 에러 객체 테스트")
	err := RegisterAccount("myId", "1234")
	if err != nil {
		if errInfo, ok := err.(PasswordError); ok {
			fmt.Printf("%v, Len:%d, RequiredLen:%d\n", errInfo, errInfo.Len, errInfo.RequireLen)
		}
	} else {
		fmt.Println("회원 가입 성공")
	}

	fmt.Println()
	fmt.Println()
}

func errorWrappingTest() {
	readEq("3 a2 2")
}

func panicTest1() {
	divide(9, 3)
	divide(9, 0)
}

/*
패닉 복구는 가능하면 하지 않는 것이 좋다.
Go는 SEH(Structured Error Handling, 구조화된 에러 처리)를 지원하지 않는다. (java의 try-catch 문)
  - 성능문제 때문
  - 에러를 먹어버리는 문제 때문
*/
func panicTest2() {
	f()
	fmt.Println("패닉 복구됨")
}

func main() {
	fmt.Println("#23 에러 핸들링")
	fileReadWriteTest()
	customErrorTest()
	customErrorTest2()
	errorWrappingTest()
	// panicTest1()
	panicTest2()
}
