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
    "encoding/base64"
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
	json.NewDecoder(r.Body).Decode(&elements)
	fmt.Println(elements)
	requests := createSetRequests(elements)
	var aggRes Response

	for i, req := range(requests) {
		if len(req) > 0 {
			url := "http://" + s[i] + "/set"
			reqBody := new(bytes.Buffer)
			json.NewEncoder(reqBody).Encode(req)
			fmt.Println("URL: ", url, "Server Req: ", reqBody)

			res, _ := resty.R().
	      SetHeader("Content-Type", "application/json").
	      SetBody(reqBody).
	      Put(url)
			fmt.Println(res.String())

			var response Response
			json.Unmarshal(res.Body(), &response)
			fmt.Println(response.KeysFailed)
			mergeRes(&aggRes, response)
		}
	}

	fmt.Println(aggRes)

	setStatusCode(w, &aggRes)

	json.NewEncoder(w).Encode(aggRes)
}

func mergeRes(aggRes *Response, res Response) {
	aggRes.KeysAdded += res.KeysAdded
	aggRes.KeysFailed = append(aggRes.KeysFailed, res.KeysFailed...)
}

func setStatusCode(w http.ResponseWriter, aggRes *Response) {
	if len(aggRes.KeysFailed) > 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func createSetRequests(elements []Element) [][]Element {
	numOfServers := len(servers())
	serverKeys := initServerKeys(numOfServers)
	for _, kv := range(elements) {
		serverIndex := hash(kv.Key.Data) % numOfServers

		if kv.Key.Encoding == "binary" {
			kv.Key.Data = base64.StdEncoding.EncodeToString([]byte(kv.Key.Data))
		}

		if kv.Value.Encoding == "binary" {
			kv.Value.Data = base64.StdEncoding.EncodeToString([]byte(kv.Value.Data))
		}

		serverKeys[serverIndex] = append(serverKeys[serverIndex], kv)
	}
	fmt.Println(serverKeys)
	return serverKeys
}

func initServerKeys(numOfServers int) [][]Element {
	serverKeys := make([][]Element, numOfServers)
	for i:=0; i < numOfServers; i++ {
		serverKeys[i] = make([]Element, 0)
	}
	return serverKeys
}

func hash(s string) int {
        h := fnv.New32a()
        h.Write([]byte(s))
        return int(h.Sum32())
}