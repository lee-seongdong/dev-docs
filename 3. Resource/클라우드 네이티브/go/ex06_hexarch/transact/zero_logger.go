package transact

import (
	"hexarch/core"
)

type zeroTransactionLogger struct{}

func (z zeroTransactionLogger) WriteDelete(key string)                        {}
func (z zeroTransactionLogger) WritePut(key, value string)                    {}
func (z zeroTransactionLogger) Err() <-chan error                             { return nil }
func (z zeroTransactionLogger) Run()                                          {}
func (z zeroTransactionLogger) Close() error                                  { return nil }
func (z zeroTransactionLogger) ReadEvents() (<-chan core.Event, <-chan error) { return nil, nil }

func NewZeroTransactionLogger() zeroTransactionLogger {
	return zeroTransactionLogger{}
}
