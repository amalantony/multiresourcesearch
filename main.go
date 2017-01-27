package main

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/amalantony/multiresourcesearch/search"
	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	query := queryParams["q"][0]
	aggregatedResults := search.Search(query)
	resultsJSON, err := json.Marshal(aggregatedResults)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s", err.Error()))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resultsJSON)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	fmt.Println("Listening on 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
