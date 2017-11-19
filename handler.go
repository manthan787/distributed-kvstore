package main

import (
    "fmt"
    "html"
    "net/http"
    "log"
    "encoding/json"
    "bytes"
)

// Type for returning and storing Query responses
// Example, {"key": {"encoding": "string", "data": "key"}, "value": true}
type QueryResponse struct {
	Key KeyValue `json:"key"`
	Value bool `json:"value"`
}

// Handles fetch requests with two possible methods:
// 1) GET /fetch => Returns all key-value pairs from all servers
// 2) POST /fetch listOfKeys => Returns all key-value pairs for given listOfKeys
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

// Handles GET /fetch requests
func fetchGetHandler(w http.ResponseWriter, r *http.Request) {
	var allElements []Element
	serverAddrs := servers()
	for _, s := range serverAddrs {
		allElements = append(allElements, fetchFromServer(s)...)
	}
    json.NewEncoder(w).Encode(allElements)
}

// Fetches all key-value pairs from given `server`
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

// Handles `POST /fetch listOfKeys` requests
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

// Groups all keys based on what server they are stored on.
// The grouped values are then used to send batch requests to each server.
func groupKeysByServer(numServers int, keys []KeyValue) [][]KeyValue{
	serverKeys := initServerKeys(numServers)
	for _, k := range(keys) {
		serverIndex := hash(k.Data) % numServers
		log.Println("ServerIndex ", serverIndex, "for key ", k.Data)
		serverKeys[serverIndex] = append(serverKeys[serverIndex], k)
	}
	return serverKeys
}

// Similar to fetchFromServer function, but only returns key-value pairs
// for given list of keys
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

// Handles /query POST requests
func queryHandler(w http.ResponseWriter, r *http.Request) {
}