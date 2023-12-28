package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func square(wg *sync.WaitGroup, ch chan int) {
	n := <-ch

	time.Sleep(time.Second)

	fmt.Println("square:", n*n)
	wg.Done()

}

func square2(wg *sync.WaitGroup, ch chan int) {
	// 무한루프
	for n := range ch {
		fmt.Println("Square: ", n*n)
		time.Sleep(time.Second)
	}
	wg.Done()
}

func square3(wg *sync.WaitGroup, ch chan int) {
	tick := time.Tick(time.Second)            // Tick : 일정 간격으로 신호를 발생하는 채널
	terminate := time.After(10 * time.Second) // After : 일정 시간 후 신호를 한번 발생하는 채널

	for {
		select {
		case <-tick:
			fmt.Println("Tick")
		case <-terminate:
			fmt.Println("Terminate")
			wg.Done()
			return
		case n, ok := <-ch:
			if ok {
				fmt.Println("Square: ", n*n)
				time.Sleep(time.Second)
			}
		}
	}
}

func channelTest() {
	fmt.Println("\n# 채널 테스트")
	var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(1)

	go square(&wg, ch)
	ch <- 9
	wg.Wait()
}

func goroutineLeakTest() {
	fmt.Println("\n# 좀비고루틴 : 채널을 닫아주지 않아서 무한대기하는 고루틴. 고루틴 릭 이라고도 한다.")
	var wg sync.WaitGroup
	ch := make(chan int)

	wg.Add(1)
	go square2(&wg, ch)

	for i := 0; i < 10; i++ {
		ch <- i * 2
	}

	close(ch) // 좀비 고루틴 방지
	wg.Wait()
}

func selectTest() {
	fmt.Println("\n# select : 여러 채널에서 동시에 데이터를 기다릴 때 사용")
	var wg sync.WaitGroup
	ch := make(chan int)

	wg.Add(1)
	go square3(&wg, ch)

	for i := 0; i < 10; i++ {
		ch <- i * 2
	}

	close(ch) // 좀비 고루틴 방지
	wg.Wait()
}

type Car struct {
	Body  string
	Tire  string
	Color string
}

var startTime = time.Now()

func MakeBody(wg *sync.WaitGroup, tireCh chan *Car) {
	tick := time.Tick(time.Second)
	after := time.After(10 * time.Second)

	for {
		select {
		case <-tick:
			car := &Car{}
			car.Body = "Sports car"
			tireCh <- car
		case <-after:
			close(tireCh)
			wg.Done()
			return
		}
	}
}

func InstallTire(wg *sync.WaitGroup, tireCh, paintCh chan *Car) {
	for car := range tireCh {
		time.Sleep(time.Second)
		car.Tire = "Winter tire"
		paintCh <- car
	}

	close(paintCh)
	wg.Done()
}

func PaintCar(wg *sync.WaitGroup, paintCh chan *Car) {
	for car := range paintCh {
		time.Sleep(time.Second)
		car.Color = "Red"
		dutation := time.Now().Sub(startTime)
		fmt.Printf("%.2f Complete Car : %s %s %s\n", dutation.Seconds(), car.Body, car.Tire, car.Color)
	}

	wg.Done()
}

func producerConsumerPatternTest() {
	fmt.Println("\n# 생산자 소비자 패턴")
	var wg sync.WaitGroup

	tireCh := make(chan *Car)
	paintCh := make(chan *Car)

	fmt.Println("Start Factory")

	wg.Add(3)

	go MakeBody(&wg, tireCh)
	go InstallTire(&wg, tireCh, paintCh)
	go PaintCar(&wg, paintCh)

	wg.Wait()
	fmt.Println("Close the Factory")
}

func PrintEverySecont(wg *sync.WaitGroup, ctx context.Context) {
	tick := time.Tick(time.Second)
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case <-tick:
			fmt.Println("tick", ctx.Value("number"), ctx.Value("text"))
		}
	}
}

func contextTest() {
	fmt.Println("\n# 컨텍스트 테스트")
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background()) // Background : 기본 컨텍스트. 추가적인 기능들은 기본 컨텍스트에 데코레이션 형식으로 사용
	ctx = context.WithValue(ctx, "number", 9)               // key, value 형식의 값을 전달할 수 있는 컨텍스트
	ctx = context.WithValue(ctx, "text", "test text")
	go PrintEverySecont(&wg, ctx)
	time.Sleep(5 * time.Second)
	cancel() // ctx.Done() 에 신호를 발생함

	wg.Wait()
}

type Publisher struct {
	ctx         context.Context
	subscribers []chan<- string // chan<-string : write only channel, <-chan string : read only channel
	subscribeCh chan chan<- string
	publisherCh chan string
}

type Subscriber struct {
	ctx   context.Context
	name  string
	msgCh chan string
}

func NewPublisher(ctx context.Context) *Publisher {
	return &Publisher{
		ctx:         ctx,
		subscribers: make([]chan<- string, 0),
		publisherCh: make(chan string),
		subscribeCh: make(chan chan<- string),
	}
}

func NewSubscriber(ctx context.Context, name string) *Subscriber {
	return &Subscriber{
		ctx:   ctx,
		name:  name,
		msgCh: make(chan string),
	}
}

func (p *Publisher) Subscribe(sub chan<- string) {
	p.subscribeCh <- sub
}

func (p *Publisher) Publish(msg string) {
	p.publisherCh <- msg
}

func (p *Publisher) Update(wg *sync.WaitGroup) {
	for {
		select {
		case sub := <-p.subscribeCh:
			p.subscribers = append(p.subscribers, sub)
		case msg := <-p.publisherCh:
			for _, subscriber := range p.subscribers {
				subscriber <- msg
			}
		case <-p.ctx.Done():
			wg.Done()
			return
		}
	}
}

func (s *Subscriber) Subscribe(pub *Publisher) {
	pub.Subscribe(s.msgCh)
}

func (s *Subscriber) Update(wg *sync.WaitGroup) {
	for {
		select {
		case msg := <-s.msgCh:
			fmt.Printf("%s got Message: %s\n", s.name, msg)
		case <-s.ctx.Done():
			wg.Done()
			return
		}
	}
}

func pubSubPatternTest() {
	fmt.Println("\n# 발행, 구독 패턴 테스트")
	var wg sync.WaitGroup
	wg.Add(4)
	ctx, cancel := context.WithCancel(context.Background())
	publisher := NewPublisher(ctx)
	subscriber1 := NewSubscriber(ctx, "AAA")
	subscriber2 := NewSubscriber(ctx, "BBB")

	go publisher.Update(&wg)

	subscriber1.Subscribe(publisher)
	subscriber2.Subscribe(publisher)

	go subscriber1.Update(&wg)
	go subscriber2.Update(&wg)

	go func() {
		tick := time.Tick(time.Second * 2)
		for {
			select {
			case <-tick:
				publisher.Publish("Hello Message")
			case <-ctx.Done():
				wg.Done()
				return
			}
		}
	}()

	fmt.Scanln()
	cancel()
	wg.Wait()
}

func main() {
	fmt.Println("#25 채널")
	fmt.Println("채널 : 고루틴끼리 메세지를 전달할 수 있는 Thread-safe-queue")
	fmt.Println("# 채널 생성")
	var messages chan int = make(chan int, 1) // 버퍼크기를 주어 비동기채널을 생성하고, 데드락 방지

	fmt.Println("# 채널에 데이터 insert")
	messages <- 5 // 버퍼 사이즈가 없는 경우, 채널이 비워질때까지 대기함.
	var msg = <-messages
	fmt.Println(msg)

	channelTest()

	// goroutineLeakTest()

	// selectTest()

	// producerConsumerPatternTest()

	// contextTest()

	pubSubPatternTest()
}
