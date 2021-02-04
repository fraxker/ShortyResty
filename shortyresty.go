package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
)

// Struct for reading json from incoming http request
type UrlRequest struct {
	Url string
}

// Struct for handling conversion to json for http response
type UrlResponse struct {
	Short_url string
}

// Seed to use for id generation
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// Map to store urls
// Where key is the id and value is the url
var m = make(map[string]string)

// Charset to use when generating ids
const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Length to use when generating id
const id_length = 8

// Port to host the server on
const server_port = "8080"

// handleShorten handles http requests for the /shorten endpoint.
// this function takes in a url from the http request,
// generates a random id for that url and adds it to the map.
// If the request is not a POST, the json is invalid or the url is invalid,
// the request will return a 400 Bad Request Error.
func handleShorten(w http.ResponseWriter, r *http.Request) {
	// Check to see if request is POST
	if r.Method != "POST" {
		// If not, throw error
		http.Error(w, "Only POST is supported at this endpoint", http.StatusBadRequest)
		return
	}
	var u UrlRequest
	// Decode URL
	err := json.NewDecoder(r.Body).Decode(&u)
	// If url cannot be decoded, throw error
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Check to see if url is valid
	_, err = url.ParseRequestURI(u.Url)
	// If not, throw erro
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	// variable to hold id for url
	var id string

	// Check to see if url is already in map
	// If so, set the url equal to that key
	for k, v := range m {
		if v == u.Url {
			id = k
			break
		}
	}
	// Otherwise, generate new id and add it to the map
	if id == "" {
		id = StringWithCharset(id_length, charset)
		m[id] = u.Url
	}

	// Assemble url to send back
	short_url := "http://127.0.0.1:" + server_port + "/" + id
	// Convert to json
	json_url := UrlResponse{short_url}
	js, err := json.Marshal(json_url)
	// Send it
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// forwardShort handles http requests from the /{id} endpoint.
// This endpoint, when given a valid id will 302 forward the request
// onto the original url from the map.
// The request will return 400 Bad Request if the ID given is not present.
func forwardShort(w http.ResponseWriter, r *http.Request) {
	// Get the id from the REST api
	vars := mux.Vars(r)
	id := vars["id"]
	// Get url from map
	url, prs := m[id]
	// If id is not present in map, throw error
	if !prs {
		http.Error(w, "ID not present", http.StatusBadRequest)
		return
	}
	// Redirect user
	http.Redirect(w, r, url, http.StatusFound)
	return
}

// StringWithCharset generates an random id of given length and given charset.
// This function is used to generate the random ids used in the handleShorten function.
// Found on https://www.calhoun.io/creating-random-strings-in-go/
func StringWithCharset(length int, charset string) string {
	// Set up byte string of length
	b := make([]byte, length)
	// for each character in length, generate random id using seededRand
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	// Convert to string and return
	return string(b)
}

// Main function of the program
// Sets up a new mux router to handle http requests for the /shorten and /{id} enpoints
func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/shorten", handleShorten)
	myRouter.HandleFunc("/{id}", forwardShort)
	log.Fatal(http.ListenAndServe(":"+server_port, myRouter))
}
