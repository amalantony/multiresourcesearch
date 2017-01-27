package search

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const (
	apiKey               = "AIzaSyD5bMRAkJmNSgE-TtsW21uTZMDMMNPE0rM"
	customSearchEngineID = "015968340880156748630%3Aui5gmgfozyu"
)

func getGoogleResults(query string, r chan Results, e chan error) {
	url := "https://www.googleapis.com/customsearch/v1?key=" + apiKey + "&cx=" + customSearchEngineID + "&q=" + query
	var myClient = &http.Client{Timeout: 1 * time.Second}
	type googleResponse struct {
		Items []struct {
			Link    string `json:"link"`
			Snippet string `json:"snippet"`
		} `json:"items"`
	}
	res, err := myClient.Get(url)
	if err != nil {
		e <- &SearchError{source: "google", message: "Request Timed out"}
		return
	}
	decoder := json.NewDecoder(res.Body)
	var response googleResponse
	for {
		if err := decoder.Decode(&response); err == io.EOF {
			break
		} else if err != nil {
			e <- &SearchError{source: "google", message: "Error parsing JSON"}
			return
		}
	}
	var resultArray []SearchResult
	for _, element := range response.Items {
		el := SearchResult{Text: element.Snippet, URL: element.Link}
		resultArray = append(resultArray, el)
	}
	googleResults := Results{source: "google", results: resultArray}
	r <- googleResults
}
