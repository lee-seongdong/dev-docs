package step1

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var transactionLogger TransactionLogger

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func helloGoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello net/http!\\n"))
}

func helloMuxHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello gorilla/mux\\n"))
}

func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	transactionLogger.WritePut(key, string(value))
	log.Printf("PUT key=%s\n", key)
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := Get(key)
	if errors.Is(err, ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value))
	log.Printf("GET key=%s\n", key)
}

func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transactionLogger.WriteDelete(key)
	log.Printf("DELETE key=%s\n", key)
}

func initializeTransactionLog() error {
	var err error

	transactionLogger, err = NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := transactionLogger.ReadEvents()
	count, e, ok := 0, Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				log.Printf("DELETE %s\n", e.Key)
				err = Delete(e.Key)
				count++
			case EventPut:
				log.Printf("PUT %s : %s\n", e.Key, e.Value)
				err = Put(e.Key, e.Value)
				count++
			}
		}
	}

	log.Printf("%d events replayed\n", count)
	transactionLogger.Run()
	return err
}

func Step1() {
	// transaction logger
	err := initializeTransactionLog()
	if err != nil {
		panic(err)
	}

	// defualtServerMux
	// http.HandleFunc("/", helloGoHandler)
	// log.Fatal(http.ListenAndServe(":8083", nil))

	// gorillaMux
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/", helloMuxHandler)
	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}
