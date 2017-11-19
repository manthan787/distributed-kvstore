package main

import (
		"os"
		"bytes"
    "fmt"
    "html"
    "log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "hash/fnv"
    "gopkg.in/resty.v1"
)

type KeyValue struct {
	Encoding string `json:"encoding"`
	Data string `json:"data"`
}

type Element struct {
	Key KeyValue `json:"key"`
	Value KeyValue `json:"value"`
}

type Response struct {
	KeysAdded int `json:"keys_added"`
	KeysFailed []KeyValue `json:"keys_failed"`
}

func main() {
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", rootHandler)
    router.HandleFunc("/set", setHandler).Methods("PUT")
    router.HandleFunc("/fetch", fetchHandler).Methods("GET", "POST")
    router.HandleFunc("/query", queryHandler).Methods("POST")
    log.Fatal(http.ListenAndServe(":8080", router))
}

func servers() []string {
	return os.Args[1:]
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(len(servers()), " servers")
  fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	s := servers()
	elements := make([]Element, 0)
	fmt.Println(r)
	json.NewDecoder(r.Body).Decode(&elements)
	requests := createSetRequests(elements)
	// var response Response

	for i := 0; i < len(requests); i++ {
		if len(requests[i]) > 0 {
			url := "http://" + s[i] + "/set"
			reqBody := new(bytes.Buffer)
			json.NewEncoder(reqBody).Encode(requests[i])
			fmt.Println("URL: ", url, "Server Req: ", reqBody)
			res, _ := resty.R().
	      SetHeader("Content-Type", "application/json").
	      SetBody(reqBody).
	      Put(url)
			fmt.Println(res.String())
			// var response Response
			// content := json.NewDecoder(res.Body()).Decode(&response)
			// fmt.Println(content)
			// response.KeysAdded += content.KeysAdded
		}
	}
}

func createSetRequests(elements []Element) [][]Element {
	numOfServers := len(servers())
	serverRequests := make([][]Element, numOfServers)
	for i := 0; i < numOfServers; i++ {
		serverRequests[i] = make([]Element, 0)
	}
	for i := 0; i < len(elements); i++ {
		kv := elements[i]
		serverIndex := hash(kv.Key.Data) % numOfServers
		serverRequests[serverIndex] = append(serverRequests[serverIndex], kv)
	}
	return serverRequests
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func hash(s string) int {
        h := fnv.New32a()
        h.Write([]byte(s))
        return int(h.Sum32())
}