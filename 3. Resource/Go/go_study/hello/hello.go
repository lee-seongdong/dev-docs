package main

import (
	"fmt"

	"study.com/greetings"
)

func main() {
	message := greetings.Hello("LSD")
	fmt.Println(message)
}
