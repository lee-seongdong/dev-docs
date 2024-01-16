package main

import (
	"hexarch/core"
	"hexarch/frontend"
	"hexarch/transact"
	"log"
)

func main() {
	tl, err := transact.NewTransactionLogger("postgres")
	if err != nil {
		log.Fatal(err)
	}

	store := core.NewKeyValueStore[string](tl)
	store.Restore()

	fe, err := frontend.NewFrontEnd("rest")
	if err != nil {
		log.Fatal(err)
	}

	err = fe.Start(store)
	if err != nil {
		log.Fatal(err)
	}
}
