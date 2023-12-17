package main

import (
	"fmt"
	"sort"
)

func changeArray(inputArray [5]int) {
	inputArray[2] = 200
}

func changeSlice(inputSlice []int) {
	inputSlice[2] = 200
}

func badAddNum(slice []int) {
	slice = append(slice, 4)
}

func goodAddNum1(slice *[]int) {
	*slice = append(*slice, 4)
}

// 1보다 조금 더 go 다운 코드
func goodAddnum2(slice []int) []int {
	return append(slice, 4)
}

type Student struct {
	name string
	age  int
}

type Students []Student

func (s Students) Len() int           { return len(s) }
func (s Students) Less(i, j int) bool { return s[i].age < s[j].age }
func (s Students) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func main() {
	fmt.Println("#18 슬라이스")
	// Slice : 동적 배열 타입 -> 가변 길이의 배열을 가리키는 포인터타입

	fmt.Println("# slice 초기화 방법")
	var slice1 []int
	slice2 := []int{1, 3, 5}
	slice3 := []int{1, 5: 10, 10: 200}
	slice4 := make([]int, 10)    // 사이즈 10인 슬라이스 초기화
	slice4 = make([]int, 10, 20) // 사이즈10, 최대용량20인 슬라이스 초기화
	arr := [...]int{1, 3, 4}     // 이것은 배열이다

	// fmt.Println(slice[1]) // panic
	fmt.Println("slice1 is empty ? ", len(slice1) == 0)
	fmt.Println("slice2 :", slice2)
	fmt.Println("slice3 :", slice3)
	fmt.Println("slice4 :", slice4)
	fmt.Println("arr :", arr)

	fmt.Println("\n# slice append")
	slice5 := append(slice4, 10)
	fmt.Println("slice5 :", slice5)
	slice5 = append(slice4, 1, 2, 3, 4, 5)
	fmt.Println("slice5 :", slice5)

	fmt.Println("\n# array와 slice의 차이")
	array2 := [5]int{1, 2, 3, 4, 5}
	slice6 := []int{1, 2, 3, 4, 5}

	changeArray(array2) // 배열을 복사해서 인자로 넘기기 때문에 array는 불변
	changeSlice(slice6) // slice는 변함. 배열값 자체를 인자로 넘기는것이 아니라 슬라이스 구조체를 복사해서 넘기기 때문에, 내부적으로 가리키는 배열은 같은 배열이다.

	fmt.Println(array2)
	fmt.Println(slice6)

	fmt.Println("\n# append() 동작원리")
	// 빈 공간(cap - len)이 충분한 경우 : 기존의 배열 포인터를 재사용
	// 빈 공간이 모자란 경우 : 새로운 배열포인터 생성, 새로운 slice를 생성
	slice7 := make([]int, 3, 5)
	appendedSlice1 := append(slice7, 4, 5)
	appendedSlice2 := append(slice7, 4, 5, 6, 7)
	fmt.Println("slice7의 cap : 5이므로 apended1은 기존의 포인터를 재사용한다.")
	fmt.Println("따라서, appended1의 특정 인덱스의 값을 변경하면, slice7에도 반영된다.")
	fmt.Println("반면에, appended2를 생성할때에는, 빈공간이 모자라 새로운 포인터를 생성했기 때문에, appended2의 특정 인덱스 값을 변경해도 slice7은 영향을 받지 않는다.")
	fmt.Printf("slice7 : %v, appended1 : %v, appended2 : %v\n", slice7, appendedSlice1, appendedSlice2)
	appendedSlice1[1] = 100
	appendedSlice2[2] = 200
	fmt.Printf("slice7 : %v, appended1 : %v, appended2 : %v\n", slice7, appendedSlice1, appendedSlice2)
	slice7[0] = 300
	fmt.Printf("slice7 : %v, appended1 : %v, appended2 : %v\n", slice7, appendedSlice1, appendedSlice2)

	fmt.Println("\n# 흔히 하는 실수")
	slice8 := []int{1, 2, 3}
	badAddNum(slice8)
	fmt.Println(slice8)
	goodAddNum1(&slice8)
	fmt.Println(slice8)
	fmt.Println(goodAddnum2(slice8))

	fmt.Println("\n# 슬라이스와 슬라이싱")
	fmt.Println("슬라이스 : 배열을 슬라이싱한 결과를 참조하는 구조체")
	arr3 := [5]int{1, 2, 3, 4, 5}
	slice9 := arr3[1:2]
	fmt.Println("array:", arr3)
	fmt.Println("slice:", slice9, len(slice9), cap(slice9))
	arr3[1] = 100
	fmt.Println("change arr3[1]")
	fmt.Println("array:", arr3)
	fmt.Println("slice:", slice9, len(slice9), cap(slice9))
	slice9 = append(slice9, 200)
	fmt.Println("append slice9 200 : 포인터는 arr3을 가리키기 때문에, arr[2] 값이 영향을 받는다")
	fmt.Println("array:", arr3)
	fmt.Println("slice:", slice9, len(slice9), cap(slice9))

	fmt.Println("\n# 슬라이스는 len이 아니라, cap을 보고 슬라이스여부를 판단한다")
	arr4 := [100]int{1: 1, 2: 2, 99: 99}
	slice10 := arr4[1:10]
	slice11 := slice10[2:99]
	fmt.Println(len(slice10), cap(slice10), slice10)
	fmt.Println(len(slice11), cap(slice11), slice11)

	fmt.Println("\n# slice[시작인덱스 : 끝인덱스 : capSize]")
	fmt.Println("capSize가 없으면, slice의 capSize만큼 슬라이스한다")

	fmt.Println("\n# slice복사 (별도의 슬라이스로 동작하도록)")
	slice12 := []int{1, 2, 3, 4, 5}
	slice13 := make([]int, len(slice12))
	for i, v := range slice12 {
		slice13[i] = v
	}
	slice13[1] = 100

	slice14 := append([]int{}, slice12...)
	slice14[1] = 200

	slice15 := make([]int, len(slice12))
	copy(slice15, slice12)
	slice15[1] = 300

	fmt.Println("origin:", slice12, "순회하여 복사:", slice13, "append를 활용하여 복사:", slice14, "copy로 복사:", slice15)

	fmt.Println("\n# 슬라이스 요소 삭제")
	index := 2
	slice16 := []int{1, 2, 3, 4, 5}
	removedSlice := append(slice16[:index], slice16[index+1:]...) // 참고. 기존 배열포인터를 참조하므로 slice16이 영향을 받는다
	fmt.Println("origin:", slice16, "removed:", removedSlice)

	fmt.Println("\n# 슬라이스 요소 삽입")
	index = 2
	slice17 := []int{1, 2, 3, 4, 5}
	insertedSlice := append(slice17[:index], append([]int{100}, slice17[index:]...)...) // 불필요한 메모리버퍼를 사용해야한다.
	fmt.Println("origin:", slice17, "insertedSlice:", insertedSlice)

	copyInsertedSlice := append(slice17, 0)
	copy(copyInsertedSlice[index+1:], copyInsertedSlice[index:])
	copyInsertedSlice[index] = 100
	fmt.Println("origin:", slice17, "copyInsertedSlice:", copyInsertedSlice)

	fmt.Println("\n# 슬라이스 정렬(primitive)")
	slice18 := []int{2, 5, 3, 67, 1, 3, 2}
	sort.Ints(slice18)
	fmt.Println(slice18)

	fmt.Println("\n# 슬라이스 정렬(structure)")
	students := []Student{
		{"화랑", 31},
		{"송하나", 52},
		{"류", 15},
		{"켄", 11},
	}

	fmt.Println("origin:", students)
	// sort.Sort(students) // 필수 인터페이스가 없기 때문에 에러
	sort.Sort(Students(students))
	fmt.Println("sorted:", students)

}
