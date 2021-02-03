package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type UrlRequest struct {
	Url string
}

type UrlResponse struct {
	Short_url string
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

var m = make(map[string]string)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const id_length = 8

const server_port = "8080"

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST is supported at this endpoint", http.StatusBadRequest)
		return
	}
	var u UrlRequest
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Generate ID
	id := StringWithCharset(id_length, charset)
	m[id] = u.Url
	short_url := "http://127.0.0.1:" + server_port + "/" + id
	json_url := UrlResponse{short_url}
	js, err := json.Marshal(json_url)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func forwardShort(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	//Get url from map
	url := m[id]
	http.Redirect(w, r, url, http.StatusFound)
	return
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/shorten", handleShorten)
	myRouter.HandleFunc("/{id}", forwardShort)
	log.Fatal(http.ListenAndServe(":"+server_port, myRouter))
}
