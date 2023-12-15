package exinit

// custompkg를 임포트 하는 경우 import cycle이 생기기 때문에 오류 발생
import "fmt"

var (
	a = c + b // 1. a 초기화를 위해 c와 b가 필요
	b = f()   // 5. b 초기화를 위해 f() 실행, 7. b는 5로 초기화
	c = f()   // 2. c 초기화를 위해 f() 실행, 4. c는 4로 초기화
	d = 3
)

// 프로그램이 초기화 될 때 한번만 호출된다. import 횟수와 상관 없이, 최초 import 되는 순간만 호출
// init() 을 통해서 패키지 내 전역변수를 초기화 할 수 있다.
// 8. init() 호출
func init() {
	d++ // 9. d = 6으로 변경
	fmt.Println("exinit 패키지 초기화", d)
}

func f() int {
	d++ // 3. d = 4로 변경, 6. d = 5로 변경
	fmt.Println("f() d:", d)
	return d
}

func PrintD() {
	fmt.Println("d:", d)
}
