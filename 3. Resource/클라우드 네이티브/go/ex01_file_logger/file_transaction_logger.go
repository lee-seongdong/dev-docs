package file_logger

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

/*
- 파일기반 로그의 장점
  - 다운스트림에 대한 의존성이 없음
  - 직관적

- 파일기반 로그의 단점
  - 확장하기 어려움
  - 파일 크기의 제한
*/
type FileTransactionLogger struct {
	events       chan<- Event
	errors       <-chan error
	lastSequence uint64
	file         *os.File
	wg           *sync.WaitGroup
}

func (l *FileTransactionLogger) WritePut(key, value string) {
	l.wg.Add(1)
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func (l *FileTransactionLogger) WriteDelete(key string) {
	l.wg.Add(1)
	l.events <- Event{EventType: EventDelete, Key: key}
}

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan Event)
	outError := make(chan error)

	go func() {
		var e Event

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			if _, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s\n", &e.Sequence, &e.EventType, &e.Key, &e.Value); err != nil {
				outError <- fmt.Errorf("input parse error: %w", err)
				return
			}

			if l.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			l.lastSequence = e.Sequence
			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()

	return outEvent, outError
}

func (l *FileTransactionLogger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		for e := range events {
			l.lastSequence++

			_, err := fmt.Fprintf(l.file, "%d\t%d\t%s\t%s\n", l.lastSequence, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
			}

			l.wg.Done()
		}
	}()
}

func (l *FileTransactionLogger) Wait() {
	l.wg.Wait()
}

func (l *FileTransactionLogger) Close() error {
	l.wg.Wait()

	if l.events != nil {
		close(l.events)
	}

	return l.file.Close()
}

func NewFileTransactionLogger(fileName string) (*FileTransactionLogger, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0775)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}

	return &FileTransactionLogger{file: file, wg: &sync.WaitGroup{}}, nil
}
