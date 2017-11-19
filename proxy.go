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

func rootHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	var aggRes Response
	s := servers()
	requests := createSetRequests(decodedReq(r))
	for i, req := range(requests) {
		if len(req) > 0 {
			var response Response
			reqBody := createReqBody(req)
			res := put("http://" + s[i] + "/set", reqBody)
			json.Unmarshal(res.Body(), &response)
			mergeRes(&aggRes, &response)
		}
	}
	setStatusCode(w, &aggRes)
	json.NewEncoder(w).Encode(aggRes)
}

func servers() []string {
	return os.Args[1:]
}

func decodedReq(r *http.Request) []Element {
	elements := make([]Element, 0)
	json.NewDecoder(r.Body).Decode(&elements)
	return elements
}

func createSetRequests(elements []Element) [][]Element {
	numOfServers := len(servers())
	serverKeys := initServerKVs(numOfServers)
	for _, kv := range(elements) {
		serverIndex := hash(kv.Key.Data) % numOfServers
		encode(&kv)
		serverKeys[serverIndex] = append(serverKeys[serverIndex], kv)
	}
	return serverKeys
}

func createReqBody(req []Element) *bytes.Buffer {
	reqBody := new(bytes.Buffer)
	json.NewEncoder(reqBody).Encode(req)
	return reqBody
}

func put(url string, reqBody *bytes.Buffer) *resty.Response {
	res, _ := resty.R().
	  SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		Put(url)
	return res
}

func mergeRes(aggRes *Response, res *Response) {
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

func initServerKVs(numOfServers int) [][]Element {
	serverKeys := make([][]Element, numOfServers)
	for i:=0; i < numOfServers; i++ {
		serverKeys[i] = make([]Element, 0)
	}
	return serverKeys
}

func encode(kv *Element) {
	if kv.Key.Encoding == "binary" {
			kv.Key.Data = base64.StdEncoding.EncodeToString([]byte(kv.Key.Data))
	}
	if kv.Value.Encoding == "binary" {
		kv.Value.Data = base64.StdEncoding.EncodeToString([]byte(kv.Value.Data))
	}
}

func hash(s string) int {
  h := fnv.New32a()
  h.Write([]byte(s))
  return int(h.Sum32())
}