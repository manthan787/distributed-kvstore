package main

import (
		"os"
    "fmt"
    "html"
    "log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
)

type KeyValue struct {
	Encoding string `json:"encoding"`
	Data string `json:"data"`
}

func main() {
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", rootHandler)
    router.HandleFunc("/set", setHandler).Methods("PUT")
    router.HandleFunc("/fetch", fetchHandler).Methods("GET", "POST")
    router.HandleFunc("/query", queryHandler).Methods("POST")
    log.Fatal(http.ListenAndServe(":8080", router))
}

func numOfServers() int {
	return len(os.Args[1:])
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	servers := os.Args[1:]
	fmt.Println(len(servers), " servers")

  fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	elements := make([]KeyValue, 0)
	json.NewDecoder(r.Body).Decode(&elements)
}

func createSetRequests(elements []KeyValue) {
	// numOfServers := numOfServers()
	// var serverRequests [][]KeyValue{}
	
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}