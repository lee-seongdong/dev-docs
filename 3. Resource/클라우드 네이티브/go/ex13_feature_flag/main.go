package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

const useNewStorage bool = false

func FeatureEnabled(flag string, r *http.Request) bool {
	return viper.GetBool(flag)
}

func NewGet(key string) (string, error) {
	return "new", nil
}

func Get(key string) (string, error) {
	return "", nil
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	var value string
	var err error
	if FeatureEnabled("use-new-storage", r) {
		value, err = NewGet(key)
	} else {
		value, err = Get(key)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value))
}

func main() {

}
