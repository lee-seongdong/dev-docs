package main

import (
	"go_study/ex_interface/fedex"
	"go_study/ex_interface/koreanPost"
)

// Duck typing 방식을 따르기 때문에, 인터페이스는 사용하는곳에서 선언한다.
// 따라서, 타입 선언 시 인터페이스 선언을 명시하지 않아도 된다.
// 덕타이핑 : 어떤 새가 오리처럼날고, 오리처럼걷고, 오리소리를 낸다면, 이 새를 오리로 취급하겠다.
// 장점 :
// - 사용자 중심의 코딩이 가능. 제공자는 구현체만 제공하고, 사용자가 필요에 따라 인터페이스를 정의해서 사용할 수 있다.
// 단점 :
// - 런타임 타입에러가 발생 할 수 있다.(go는 컴파일 언어이므로, 해당없음)
// - 객체 구현 시 인터페이스를 고려하려 개발해야 하기 때문에, 개발난이도가 올라간다.
type Sender interface {
	Send(parcel string)
}

//	func SendBook(name string, sender *fedex.FedexSender) {
//		sender.Send(name)
//	}

//	func SendBook(name string, sender *koreanPost.PostSender) {
//		sender.Send(name)
//	}

func SendBook(name string, sender Sender) {
	sender.Send(name)
}

func main() {
	// sender := &fedex.FedexSender{}
	var fedexSender Sender = &fedex.FedexSender{}
	var postSender Sender = &koreanPost.PostSender{}

	SendBook("어린왕자", fedexSender)
	SendBook("그리스인 조르바", fedexSender)
	SendBook("어린왕자", postSender)
	SendBook("그리스인 조르바", postSender)
}
