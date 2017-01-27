package search

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func getDuckDuckGoResults(query string, r chan Results, e chan error) {
	url := "http://api.duckduckgo.com/?q=" + query + "&format=json"
	var myClient = &http.Client{Timeout: 1 * time.Second}
	type duckduckgoResponse struct {
		RelatedTopics []struct {
			FirstURL string `json:"FirstURL"`
			Text     string `json:"Text"`
		} `json:"RelatedTopics"`
	}
	res, err := myClient.Get(url)
	if err != nil {
		e <- &searchError{source: "duckduckgo", message: "Request Timed out"}
		return
	}
	decoder := json.NewDecoder(res.Body)
	var response duckduckgoResponse
	for {
		if err := decoder.Decode(&response); err == io.EOF {
			break
		} else if err != nil {
			e <- &searchError{source: "duckduckgo", message: "Error parsing JSON"}
			return
		}
	}
	resultArray := make([]searchResult, len(response.RelatedTopics))
	for _, element := range response.RelatedTopics {
		el := searchResult{text: element.Text, url: element.FirstURL}
		resultArray = append(resultArray, el)
	}
	duckduckgoResults := Results{source: "duckduckgo", results: resultArray}
	r <- duckduckgoResults
}
