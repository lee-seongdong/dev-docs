package main

import (
	"fmt"
)

// iota 를 사용한 열거형
const (
	Red = iota
	Blue
	Green
)

const (
	BitFlag1 uint = 1 << iota
	BitFlag2
	BitFlag3
	BitFlag4
)

const (
	MasterRoom uint8 = 1 << iota
	LivingRoom
	BathRoom
	SmallRoom
)

func turnOn(rooms, room uint8) uint8 {
	return rooms | room
}

func turnOff(rooms, room uint8) uint8 {
	return rooms &^ room // 비트 클리어 연산. ^ 연산으로 room 비트를 모두 반전 후 rooms와 & 연산
}

func isTurnOn(rooms, room uint8) bool {
	return rooms&room == room
}

func l(rooms uint8) {
	if isTurnOn(rooms, MasterRoom) {
		fmt.Println("master is turn on")
	}
	if isTurnOn(rooms, LivingRoom) {
		fmt.Println("living is turn on ")
	}
	if isTurnOn(rooms, BathRoom) {
		fmt.Println("bath is turn on")
	}
	if isTurnOn(rooms, SmallRoom) {
		fmt.Println("small is turn on")
	}
}

func main() {
	fmt.Println(Red, Blue, Green)
	fmt.Println(BitFlag1, BitFlag2, BitFlag3, BitFlag4)

	var rooms uint8 = 0

	rooms = turnOn(rooms, MasterRoom)
	rooms = turnOn(rooms, BathRoom)
	rooms = turnOn(rooms, SmallRoom)

	rooms = turnOff(rooms, BathRoom)

	l(rooms)

	const PI = 3.14
	const fPI float64 = 3.14

	var a int = PI * 100
	// var b int = fPI * 100  // 에러!
	fmt.Println(a)
}
