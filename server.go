package main

import (
	"os"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/valep27/shortener/data"
	"github.com/valep27/shortener/transform"
)

var shortener *transform.UrlShortener
var store data.Storage

// Config is used to configure the shortener and its backing storage.
type Config struct {
	Alphabet      string `json:"alphabet,omitempty"`
	Salt          string `json:"salt,omitempty"`
	RedisAddress  string `json:"redisAddress,omitempty"`
	RedisPassword string `json:"redisPassword,omitempty"`
	RedisDb       int    `json:"redisDb,omitempty"`
}

func loadConfiguration() Config {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Could not find configuration file")
	}

	decoder := json.NewDecoder(configFile)
	conf := Config {}

	err = decoder.Decode(&conf)
	if err != nil {
		log.Fatalf("Could not parse configuration file: %s", err.Error())
	}

	return conf
}

func main() {
	config := loadConfiguration()

	if len(config.Alphabet) != 0 {
		shortener = transform.NewShortenerWithAlphabet(config.Salt, config.Alphabet)
	} else {
		shortener = transform.NewShortener(config.Salt)
	}

	store = data.NewRedisStorage(config.RedisAddress, config.RedisPassword, config.RedisDb)

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

	url, err := store.Get(key)
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

	index, err := store.Next()
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

	err = store.Set(id, req.URL)
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

	url, err := store.Get(key)
	if err != nil {
		log.Printf("error when getting key <%s>: %s", key, err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, url)
}
