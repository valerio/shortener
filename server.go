package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/valep27/shortener/data"
	"github.com/valep27/shortener/transform"
)

var shortener = transform.NewShortener("salt")

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/urls", postURLHandler).Methods("POST")
	r.HandleFunc("/urls/{key}", getURLHandler).Methods("GET")
	
	r.HandleFunc("/{key}", redirectHandler)
	
	log.Println("Starting url shortener service on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	url, err := data.Store.Get(key)
	if err != nil {
		log.Printf("redirect - key <%s> not found: %s", key, err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func postURLHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type request struct {
		URL string `json:"url"`
	}
	var req request

	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("error when decoding post url: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	index, err := data.Store.Next()
	if err != nil {
		log.Printf("error when incrementing index: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := shortener.Encode(index)
	if err != nil {
		log.Printf("could not encode url: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = data.Store.Set(id, req.URL)
	if err != nil {
		log.Printf("error when storing url: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, id)
}

func getURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	url, err := data.Store.Get(key)
	if err != nil {
		log.Printf("error when getting key <%s>: %s", key, err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, url)
}
