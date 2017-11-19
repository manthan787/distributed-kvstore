package main

import (
	"os"
	"bytes"
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "hash/fnv"
    "gopkg.in/resty.v1"
    "encoding/base64"
    "path"
   	"net/url"
)

// Encoding and Data of a Key or Value
type KeyValue struct {
	Encoding string `json:"encoding"`
	Data string `json:"data"`
}

// Key and Value for key-value store
type Element struct {
	Key KeyValue `json:"key"`
	Value KeyValue `json:"value"`
}

// Response given by proxy server to client
type Response struct {
	KeysAdded int `json:"keys_added"`
	KeysFailed []KeyValue `json:"keys_failed"`
}

// Create proxy server to listen at localhost:8080
// Handle valid routes: /, set, fetch, query
func main() {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", rootHandler)
  router.HandleFunc("/set", setHandler).Methods("PUT")
  router.HandleFunc("/fetch", fetchHandler).Methods("GET", "POST")
  router.HandleFunc("/query", queryHandler).Methods("POST")
  log.Fatal(http.ListenAndServe(":8080", router))
}

// Simple root handler to see if server is working
func rootHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "KV Proxy is working.")
}

func getServerPath(server string, p string) string {
	u, _ := url.Parse("http://" + server)
	u.Path = path.Join(u.Path, p)
	return u.String()
}


// Handle set PUT request to store given key-value batch as
// [ {key: {encoding: , data:}, value: {encoding: , data:}}, ... ]
func setHandler(w http.ResponseWriter, r *http.Request) {
	var aggRes Response
	s := servers()
	requests := createSetRequests(decodedReq(r))
	for i, req := range(requests) {
		if len(req) > 0 {
			reqBody := createReqBody(req)
			response := put(getServerPath(s[i], "/set"), reqBody)
			mergeRes(&aggRes, &response)
		}
	}
	setStatusCode(w, &aggRes)
	json.NewEncoder(w).Encode(aggRes)
}

// give list of server addresses given as argument **without http://** 
func servers() []string {
	return os.Args[1:]
}

// give decoded http json request as an array of Element
func decodedReq(r *http.Request) []Element {
	elements := make([]Element, 0)
	json.NewDecoder(r.Body).Decode(&elements)
	return elements
}

// return requests grouped by server based on key
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

// return request body by encoding array of Element
func createReqBody(req []Element) *bytes.Buffer {
	reqBody := new(bytes.Buffer)
	json.NewEncoder(reqBody).Encode(req)
	return reqBody
}

// perform PUT for given url and request body and return Response
func put(url string, reqBody *bytes.Buffer) Response {
	var response Response
	res, _ := resty.R().
	  SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		Put(url)
	json.Unmarshal(res.Body(), &response)
	return response
}

// merge given response per server with aggregate response for client
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