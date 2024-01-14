package transact

import (
	"fmt"
	"hexarch/core"
)

func NewTransactionLogger(logger string) (core.TransactionLogger, error) {
	switch logger {
	case "zero":
		return NewZeroTransactionLogger(), nil
	case "file":
		return NewFileTransactionLogger("./transactions.txt")
	case "postgres":
		return NewPostgresTransactionLogger(postgresDbParams{
			host:     "localhost",
			dbName:   "lsd",
			user:     "postgres",
			password: "1234",
		})
	case "":
		return nil, fmt.Errorf("transaction logger type not defined")
	default:
		return nil, fmt.Errorf("no such transaction logger %s", logger)
	}
}
