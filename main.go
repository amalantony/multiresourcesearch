package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"encoding/json"

	"github.com/amalantony/multiresourcesearch/search"
	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	w.Header().Set("Content-Type", "application/json")
	if len(queryParams) < 1 {
		errStructure := struct {
			Message string
		}{Message: "No Query Specified"}
		errJSON, _ := json.Marshal(errStructure)
		w.Write(errJSON)
		return
	}
	query := queryParams["q"][0]
	aggregatedResults := search.Search(query)
	resultsJSON, err := json.Marshal(aggregatedResults)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s", err.Error()))
	}
	w.Write(resultsJSON)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	fmt.Println("Listening on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
