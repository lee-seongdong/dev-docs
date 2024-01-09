package file_logger

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	ReadEvents() (<-chan Event, <-chan error)
	Err() <-chan error
	Run()
}

type EventType byte

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}
