package main

import (
	"fmt"
	"go_study/study_16_usepkg/custompkg"
	"go_study/study_16_usepkg/exinit"

	"github.com/guptarohit/asciigraph"
	"github.com/tuckersGo/musthaveGo/ch16/expkg"
)

// go mod tidy로 사용하는 패키지 다운로드(GOPATH/pkg/mod : /Users/user/go/pkg/mod) 및 종속성 업데이트(go.sum 파일)

func main() {
	fmt.Println("#16 패키지 사용")
	custompkg.PrintCustom()
	expkg.PrintSample()

	data := []float64{3, 4, 5, 6, 9, 7, 5, 8, 5, 10, 2, 7, 2, 5, 6}
	graph := asciigraph.Plot(data)
	fmt.Println(graph)

	fmt.Println("# exinit")
	custompkg.PrintCustom()
	exinit.PrintD()

}
