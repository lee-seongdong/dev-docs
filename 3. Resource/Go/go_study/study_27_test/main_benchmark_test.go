// 1. 테스트코드 파일명은 _test.go 로 끝나야 한다.
// 2. testing 패키지를 임포트 해야 한다.
// 3. 벤치마크 코드는 func BenchmarkXxxx(b *testing.B) 형태여야 한다
package main

import (
	"testing"
)

func BenchmarkFibonacci1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fibonacci1(20)
	}
}

func BenchmarkFibonacci2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fibonacci2(20)
	}
}
