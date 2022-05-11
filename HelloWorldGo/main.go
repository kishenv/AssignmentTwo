// A simple HTTP program to respond "Hello Word" on accessing /helloworld endpoint
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	log.Print("HelloWorld HTTP Server")
	http.HandleFunc("/helloworld", helloWorld)
	http.ListenAndServe(":8080", nil)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Print("Endpoint /hellowold accessed")
	resp, err := json.Marshal(`{"Message":"Hello World"}`)
	if err != nil {
		log.Fatalf("Error during marshal : %s", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
