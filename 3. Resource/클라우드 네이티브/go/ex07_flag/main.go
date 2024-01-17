package main

import (
	"flag"
	"fmt"
)

func main() {
	strp := flag.String("string", "foo", "a string")
	intp := flag.Int("number", 42, "an integer")
	boolp := flag.Bool("boolean", false, "a boolean")

	flag.Parse()

	fmt.Println("string:", *strp)
	fmt.Println("integer:", *intp)
	fmt.Println("boolean:", *boolp)
	fmt.Println("args:", flag.Args())
}
