package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type FibonacciJson struct {
	Position int `json:"position"`
	Element int `json:"element"`
}

func fibonacciCalculator(w http.ResponseWriter, r *http.Request){
	// read parameter from url
	parameter := r.URL.Query().Get("a-param")

	// assure that parameter can be converted to integer
	n, err := strconv.Atoi(parameter)

	if err != nil {
    	fmt.Fprintf(w,"400 Bad Request: You must provide an integer.")
    } else {
		// calculate n'th element of fibonacci sequence
		f := make([]int, n+1, n+2)
		if n < 2 {
			f = f[0:2]
		}
		f[0], f[1] = 0,1
		for i := 2; i <= n; i++{
			f[i] = f[i-1] + f[i-2]
		}

		// create response
		jsonFile := FibonacciJson{Position:n, Element:f[n]}

		// expose response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonFile)
	}
}

func main() {
	// init router for building API
	r := mux.NewRouter()

	// route handles & endpoints
	r.HandleFunc("/", fibonacciCalculator).Methods("GET")

	// start server that listens to 8080. wrap with log for error checking
	log.Fatal(http.ListenAndServe(":8080", r))
}