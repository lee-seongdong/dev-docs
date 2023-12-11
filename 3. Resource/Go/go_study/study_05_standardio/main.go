package main

import (
	"bufio"
	"fmt"
	"os"
)

func standardOut() {
	a := 10
	b := 20
	f := 32799438743.8297

	fmt.Println("Standard Out")
	fmt.Print("Print a:", a, "b:", b)
	fmt.Println("Println a:", a, "b:", b, "f:", f)    // f가 지수표현으로 출력
	fmt.Printf("Printf a: %d b: %d f: %f\n", a, b, f) // f가 실수표현으로 출력
}

func standardIn() {
	var a, b int
	n, err := fmt.Scanln(&a, &b)

	fmt.Println("Standard In")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(n, a, b)
	}
}

func buffer() {
	stdin := bufio.NewReader(os.Stdin)
	var a, b int

	n, err := fmt.Scanln(&a, &b)
	if err != nil {
		fmt.Println(err)
		stdin.ReadString('\n')
	} else {
		fmt.Println(n, a, b)
	}

	n, err = fmt.Scanln(&a, &b)
	if err != nil {
		fmt.Println(err)
		stdin.ReadString('\n')
	} else {
		fmt.Println(n, a, b)
	}
}

func main() {
	// standardOut()
	// standardIn()
	buffer()
}
