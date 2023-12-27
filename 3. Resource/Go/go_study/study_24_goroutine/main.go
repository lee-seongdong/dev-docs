package main

import (
	"fmt"
	"sync"
	"time"
)

type Account struct {
	Balance int
}

type Job interface {
	Do()
}

type SquareJob struct {
	index int
}

func (j *SquareJob) Do() {
	fmt.Printf("%d work start\n", j.index)
	time.Sleep(1 * time.Second)
	fmt.Printf("%d work end - result %d\n", j.index, j.index*j.index)
}

func PrintHangle() {
	hangle := []rune{'가', '나', '다', '라', '마', '바'}
	for _, v := range hangle {
		time.Sleep(300 * time.Millisecond)
		fmt.Printf("%c ", v)
	}
}

func PrintNumber() {
	for i := 1; i <= 5; i++ {
		time.Sleep(400 * time.Millisecond)
		fmt.Printf("%d ", i)
	}
}

var gwg sync.WaitGroup

func SumAtoB(a, b int) {
	sum := 0
	for i := a; i <= b; i++ {
		sum += i
	}
	fmt.Printf("%d부터 %d까지 합계는 %d입니다.\n", a, b, sum)
	gwg.Done()
}

func test1() {
	go PrintHangle()
	go PrintNumber()

	// time.Sleep(3 * time.Second) // 메인 고루틴이 종료되면, 다른 서브 고루틴들도 모두 종료됨
}

func test2() {
	const numOfGoRoutine = 10

	gwg.Add(numOfGoRoutine)
	for i := 0; i < numOfGoRoutine; i++ {
		go SumAtoB(1, 1000000000)
	}

	gwg.Wait()
}

var mutex sync.Mutex

func DepositAndWithdraw(account *Account) {
	mutex.Lock()
	defer mutex.Unlock()

	if account.Balance < 0 {
		panic(fmt.Sprintf("Balance : %d", account.Balance))
	}
	account.Balance += 1000
	time.Sleep(time.Millisecond)
	account.Balance -= 1000
}

func test3() {
	var wg2 sync.WaitGroup
	account := &Account{10}

	wg2.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				DepositAndWithdraw(account)
			}
			wg2.Done()
		}()
	}
	wg2.Wait()
}

func test4() {
	var jobList [10]Job
	for i := 0; i < 10; i++ {
		jobList[i] = &SquareJob{i}
	}

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		job := jobList[i]
		go func() {
			job.Do()
			wg.Done()
		}()
	}

	wg.Wait()
}

func main() {
	fmt.Println("#24 고루틴")
	fmt.Println(`- 고루틴은 Go 프로그램에서 실행되는 경량스레드
- Go 프로그램은 최소 1개의 고루틴을 가진다 (main()을 실행하는 메인고루틴)
- OS스레드를 이용하는 경량 스레드 (고루틴 != 스레드)
- CPU 코어와 OS 스레드를 1:1로 매핑하여 생성하고, 이를 이용한다. (코어 - OS스레드 - 고루틴 으로 매핑)
- 시스템콜 호출 시 고루틴은 대기상태로 바뀌고, 대기중인 고루틴으로 교체된다.
- 코어와 OS스레드가 1:1로 매핑되어 유지되므로 OS단에서 컨텍스트 스위칭 비용이 발생하지 않는다.
  (고루틴 컨텍스트 스위칭은 발생하지만, 고루틴 컨텍스트는 OS에 비해서는 아주 경량이기 때문에 비용이 아주 적게 발생함.)
- Go 프로그램은 최소 1개의 고루틴을 가진다 (main()을 실행하는 메인고루틴)`)

	// test1()
	// test2()
	// test3()

	fmt.Println(`
- 멀티스레드 환경에서는 동시성 문제가 발생한다.
- 해결방법
  - 뮤텍스 : 동시성 프로그래밍으로 인해 성능 향상을 얻을 수 없음. 오히려 하락할 수 도 있음. 데드락 발생 가능성 있음. -> 최소한의 범위에서만 사용해야 한다.
  - 영역 분할 기법
  - 역할 분할 기법`)

	test4()
}
