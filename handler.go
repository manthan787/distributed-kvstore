package main

import (
    "fmt"
    "html"
    "net/http"
    "log"
    "encoding/json"
    "bytes"
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
	keys := make([]KeyValue, 0)
	json.NewDecoder(r.Body).Decode(&keys)
	log.Println(keys)
	numServers := len(servers())
	servs := servers()
	serverKeys := groupKeysByServer(numServers, keys)
	result := make([]Element, 0)
	for idx, keys := range(serverKeys) {
		encodedList, err := json.Marshal(keys)
		if err != nil {
			log.Println("Error marshalling list of keys:", err)
			break
		}
		log.Println(string(encodedList))
		els := fetchListFromServer(servs[idx], encodedList)
		log.Println(els)
		result = append(result, els...)
	}
	json.NewEncoder(w).Encode(result)
}

func groupKeysByServer(numServers int, keys []KeyValue) [][]KeyValue{
	serverKeys := initServerKeys(numServers)
	for _, k := range(keys) {
		serverIndex := hash(k.Data) % numServers
		log.Println("ServerIndex ", serverIndex, "for key ", k.Data)
		serverKeys[serverIndex] = append(serverKeys[serverIndex], k)
	}
	return serverKeys
}

func fetchListFromServer(server string, list []byte) []Element {
	log.Println("Querying Server ", server)
	elements := make([]Element, 0)
	resp, err := http.Post(server + "/fetch", "application/json", bytes.NewBuffer(list))
	if err != nil {
		log.Fatal("Error", err)
	}
	json.NewDecoder(resp.Body).Decode(&elements)
	return elements
}

func initServerKeys(numServers int) [][]KeyValue {
	serverKeys := make([][]KeyValue, numServers)
	for i:=0; i < numServers; i++ {
		serverKeys[i] = make([]KeyValue, 0)
	}
	return serverKeys
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}