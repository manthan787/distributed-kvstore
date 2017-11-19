package main

import (
    "fmt"
    "html"
    "net/http"
    "log"
    // "io/ioutil"
    "encoding/json"
)

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			fetchGetHandler(w, r)
		case "POST":
			fetchPostHandler(w, r)
		default:
			log.Println("Invalid Method for /fetch")
	}
}

func fetchGetHandler(w http.ResponseWriter, r *http.Request) {
	var allElements []Element
	serverAddrs := servers()
	for _, s := range serverAddrs {
		allElements = append(allElements, fetchFromServer(s)...)
	}
    json.NewEncoder(w).Encode(allElements)
}

func fetchFromServer(server string) []Element {
	log.Println("Querying Server ", server)
	elements := make([]Element, 0)
	resp, err := http.Get(server + "/fetch")
	if err != nil {
		log.Fatal("Error", err)
	}
	json.NewDecoder(resp.Body).Decode(&elements)
	return elements
}

func fetchPostHandler(w http.ResponseWriter, r *http.Request) {

}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}