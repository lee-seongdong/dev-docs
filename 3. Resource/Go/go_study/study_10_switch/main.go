package main

import (
	"fmt"
)

type ColorType int // 타입 별칭

const (
	Red ColorType = iota
	Blue
	Green
	Yellow
)

func colorToString(color ColorType) string {
	switch color {
	case Red:
		return "Red"
	case Blue:
		return "Blue"
	case Green:
		return "Green"
	default:
		return "undefined"
	}
}

func getMyFavoriteColor() ColorType {
	return Blue
}

func main() {
	// switch 초기문; 비굣값
	switch a, b := 1, 2; a + b {
	case 1:
		fmt.Println("res = 1")
		// break 없어도 빠져나간다
		// 아래쪽의 case를 실행하려면 fallthrough
	case 2:
		fmt.Println("res = 2")
	case 3:
		fmt.Println("res = 3")
	}

	fmt.Println("favorite color index is", getMyFavoriteColor())
	fmt.Println("favorite color is", colorToString(getMyFavoriteColor()))
	fmt.Println("favorite color is", colorToString(Blue))
}
