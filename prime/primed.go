package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
	"strconv"
)

type PrimeJson struct {
	Number int `json:"number"`
	Prime bool `json:"prime"`
}

func primeCalculator(w http.ResponseWriter, r *http.Request){
	// read parameter from url
	parameter := r.URL.Query().Get("a-param")

	// assure that parameter can be converted to integer
	n, err := strconv.Atoi(parameter)

	if err != nil {
    	fmt.Fprintf(w,"400 Bad Request: You must provide an integer.")
    } else {
    	// check if n is prime
		primeChecker := true
	    for i := 2; i <= int(math.Floor(math.Sqrt(float64(n)))); i++ {
	        if n%i == 0 {
	            primeChecker = false
	            break
	        }
	    }

		// create response
		jsonFile := PrimeJson{Number:n, Prime:primeChecker}

		// expose response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonFile)
	}
}

func main() {
	// init router for building API
	r := mux.NewRouter()

	// route handles & endpoints
	r.HandleFunc("/", primeCalculator).Methods("GET")

	// start server that listens to 8080. wrap with log for error checking
	log.Fatal(http.ListenAndServe(":8080", r))
}