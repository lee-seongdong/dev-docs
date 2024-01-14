package transact

import (
	"bufio"
	"fmt"
	"hexarch/core"
	"os"
	"sync"
)

type fileTransactionLogger struct {
	events       chan<- core.Event
	errors       <-chan error
	lastSequence uint64
	file         *os.File
	wg           *sync.WaitGroup
}

func (l *fileTransactionLogger) WritePut(key, value string) {
	l.wg.Add(1)
	l.events <- core.Event{EventType: core.EventPut, Key: key, Value: value}
}

func (l *fileTransactionLogger) WriteDelete(key string) {
	l.wg.Add(1)
	l.events <- core.Event{EventType: core.EventDelete, Key: key}
}

func (l *fileTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *fileTransactionLogger) Run() {
	events := make(chan core.Event, 16)
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

func (l *fileTransactionLogger) Close() error {
	l.wg.Wait()

	if l.events != nil {
		close(l.events)
	}

	return l.file.Close()
}

func (l *fileTransactionLogger) ReadEvents() (<-chan core.Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan core.Event)
	outError := make(chan error)

	go func() {
		var e core.Event

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

func NewFileTransactionLogger(fileName string) (core.TransactionLogger, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0775) // read and write, append, create if not exists
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}

	return &fileTransactionLogger{file: file, wg: &sync.WaitGroup{}}, nil
}
