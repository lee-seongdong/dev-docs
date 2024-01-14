package transact

import (
	"database/sql"
	"fmt"
	"hexarch/core"
	"net/url"
	"sync"
)

type postgresDbParams struct {
	dbName   string
	host     string
	user     string
	password string
}

type postgresTransactionLogger struct {
	events chan<- core.Event // Write-only channel for sending events
	errors <-chan error      // Read-only channel for receiving errors
	db     *sql.DB           // Our database access interface
	wg     *sync.WaitGroup   // Used to ensure writes are completed
}

func (l *postgresTransactionLogger) WritePut(key, value string) {
	l.wg.Add(1)
	l.events <- core.Event{EventType: core.EventPut, Key: key, Value: url.QueryEscape(value)}
	l.wg.Done()
}

func (l *postgresTransactionLogger) WriteDelete(key string) {
	l.wg.Add(1)
	l.events <- core.Event{EventType: core.EventDelete, Key: key}
	l.wg.Done()
}

func (l *postgresTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *postgresTransactionLogger) Run() {
	events := make(chan core.Event, 16) // Make an events channel
	l.events = events

	errors := make(chan error, 1) // Make an errors channel
	l.errors = errors

	go func() { // The INSERT query
		query := `INSERT INTO transactions
			(event_type, key, value)
			VALUES ($1, $2, $3)`

		for e := range events { // Retrieve the next Event
			_, err := l.db.Exec( // Execute the INSERT query
				query,
				e.EventType, e.Key, e.Value)

			if err != nil {
				errors <- err
			}
		}
	}()
}

func (l *postgresTransactionLogger) Close() error {
	l.wg.Wait()

	if l.events != nil {
		close(l.events) // Terminates Run loop and goroutine
	}

	return l.db.Close()
}

func (l *postgresTransactionLogger) ReadEvents() (<-chan core.Event, <-chan error) {
	outEvent := make(chan core.Event) // unbuffered events channel
	outError := make(chan error, 1)   // buffered errors channel

	query := "SELECT sequence, event_type, key, value FROM transactions"

	go func() {
		defer close(outEvent)
		defer close(outError)

		rows, err := l.db.Query(query)
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}

		defer rows.Close() // This is important!

		var e core.Event

		for rows.Next() {
			err = rows.Scan(&e.Sequence, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				outError <- err
				return
			}

			outEvent <- e
		}

		err = rows.Err()
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}

func (l *postgresTransactionLogger) verifyTableExists() (bool, error) {
	const table = "transactions"

	var result string

	rows, err := l.db.Query(fmt.Sprintf("SELECT to_regclass('public.%s');", table))
	defer rows.Close()
	if err != nil {
		return false, err
	}

	for rows.Next() && result != table {
		rows.Scan(&result)
	}

	return result == table, rows.Err()
}

func (l *postgresTransactionLogger) createTable() error {
	var err error

	createQuery := `CREATE TABLE transactions (
		sequence      BIGSERIAL PRIMARY KEY,
		event_type    SMALLINT,
		key 		  TEXT,
		value         TEXT
	  );`

	_, err = l.db.Exec(createQuery)
	if err != nil {
		return err
	}

	return nil
}

func NewPostgresTransactionLogger(param postgresDbParams) (core.TransactionLogger, error) {
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", param.host, param.dbName, param.user, param.password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	logger := &postgresTransactionLogger{db: db}

	exists, err := logger.verifyTableExists()
	if err != nil {
		return nil, fmt.Errorf("failed to verify table exists: %w", err)
	}

	if !exists {
		if err = logger.createTable(); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	return logger, nil
}
