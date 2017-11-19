package main

import (
    "net/http"
    "log"
    "encoding/json"
    "bytes"
    "io"
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
	log.Println(getServerPath(server, "/fetch"))
	elements := make([]Element, 0)
	resp, err := http.Get(getServerPath(server, "/fetch"))
	if err != nil {
		log.Fatal("Error", err)
	}
	json.NewDecoder(resp.Body).Decode(&elements)
	return elements
}

// Handles `POST /fetch listOfKeys` requests
func fetchPostHandler(w http.ResponseWriter, r *http.Request) {
	keys := readKeys(r.Body)
	log.Println(keys)
	servs := servers()
	numServers := len(servs)
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
	if len(keys) > len(result) {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(result)
}

func readKeys(body io.ReadCloser) []KeyValue {
	keys := make([]KeyValue, 0)
	err := json.NewDecoder(body).Decode(&keys)
	if err != nil {
		log.Fatal("Error decoding json ", err)
	}
	return keys
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
	log.Println("URL ", getServerPath(server, "/fetch"))
	elements := make([]Element, 0)
	resp, err := http.Post(getServerPath(server, "/fetch"), "application/json", bytes.NewBuffer(list))
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
	keys := readKeys(r.Body)
	servs := servers()
	numServers := len(servs)
	serverKeys := groupKeysByServer(numServers, keys)
	result := make([]QueryResponse, 0)
	for idx, keys := range(serverKeys) {
		encodedList, err := json.Marshal(keys)
		if err != nil {
			log.Println("Error marshalling list of keys:", err)
			break
		}
		log.Println(string(encodedList))
		els := fetchQueryRespFromServer(servs[idx], encodedList)
		log.Println(els)
		result = append(result, els...)
	}
	if len(keys) > len(result) {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(result)
}

// Similar to fetchListFromServer function, but this function returns `QueryResponse`
// instead of `Element`
func fetchQueryRespFromServer(server string, list []byte) []QueryResponse {
	log.Println("Querying Server ", server)
	responses := make([]QueryResponse, 0)
	resp, err := http.Post(server + "/fetch", "application/json", bytes.NewBuffer(list))
	if err != nil {
		log.Fatal("Error", err)
	}
	json.NewDecoder(resp.Body).Decode(&responses)
	return responses
}