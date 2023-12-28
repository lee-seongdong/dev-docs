// 1. 테스트코드 파일명은 _test.go 로 끝나야 한다.
// 2. testing 패키지를 임포트 해야 한다.
// 3. 테스트 코드는 func TestXxxx(t *testing.T) 형태여야 한다
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -run TestSquare 명령어로 일부만 테스트 가능
func TestSquare(t *testing.T) {
	rst := square(9)
	if rst != 81 {
		t.Errorf("square(9) should be 81 but returns %d", rst)
	}
}

func TestSquare2(t *testing.T) {
	rst := square(3)
	if rst != 9 {
		t.Errorf("square(3) should be 9 but returns %d", rst)
	}
}

func TestSquare3(t *testing.T) {
	assert.Equal(t, square(9), 81, "square(3) should be 9")
}
