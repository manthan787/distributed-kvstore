package main

import (
	"fmt"
	"net/http"
	// "encoding/json"
	"log"
	"io/ioutil"
)

type json_data struct {
	keyValue map[string] string
}

type keyvalue_meta struct {
	Encoding string `json:"encoding"`
	Data string `json:"data"`
}

type keyvalue struct {
	Key keyvalue_meta `json:"key"`
	Value keyvalue_meta `json:"value"`
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:3000/fetch")
	if err != nil {
		log.Fatal("Error", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error While reading response: ", err)
	}
	log.Println("String", string(body))
    fmt.Fprintf(w, "%s", body)
}

func main() {
    http.HandleFunc("/fetch/", fetchHandler)
    log.Println("Staring server at", 8080)
	http.ListenAndServe(":8080", nil)
}