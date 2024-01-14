package core

import (
	"errors"
	"log"
)

var ErrorNoSuchKey = errors.New("no such key")

type KeyValueStore struct {
	m        map[string]string
	transact TransactionLogger
}

func NewKeyValueStore[T any](tl TransactionLogger) *KeyValueStore {
	return &KeyValueStore{
		m:        make(map[string]string),
		transact: tl,
	}
}

func (store *KeyValueStore) Delete(key string) error {
	delete(store.m, key)
	store.transact.WriteDelete(key)
	return nil
}

func (store *KeyValueStore) Put(key, value string) error {
	store.m[key] = value
	store.transact.WritePut(key, value)
	return nil
}

func (store *KeyValueStore) Get(key string) (string, error) {
	val, ok := store.m[key]
	if !ok {
		return val, errors.New("key not found")
	}

	return store.m[key], nil
}

func (store *KeyValueStore) WithTransactionLogger(tl TransactionLogger) *KeyValueStore {
	store.transact = tl
	return store
}

func (store *KeyValueStore) Restore() error {
	var err error

	events, errors := store.transact.ReadEvents()
	count, ok, e := 0, true, Event{}

	for ok && err == nil {
		select {
		case err, ok = <-errors:

		case e, ok = <-events:
			switch e.EventType {
			case EventDelete: // Got a DELETE event!
				err = store.Delete(e.Key)
				count++
			case EventPut: // Got a PUT event!
				err = store.Put(e.Key, e.Value)
				count++
			}
		}
	}

	log.Printf("%d events replayed\n", count)

	store.transact.Run()

	go func() {
		for err := range store.transact.Err() {
			log.Print(err)
		}
	}()

	return err
}
