package core

type EventType byte

const (
	_                     = iota // iota == 0; ignore this value
	EventDelete EventType = iota // iota == 1
	EventPut                     // iota == 2; implicitly repeat last
)

type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}

// 코어 로직이 드리븐 어댑터에 대해 동작해야 하므로 코어 애플리케이션은 포트를 알아야 한다.
// 포트 = TransactionLogger 인터페이스
type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error

	Run()
	Close() error

	ReadEvents() (<-chan Event, <-chan error)
}
