package search

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

func getDuckDuckGoResults(query string, r chan Results, e chan error) {
	Url, _ := url.Parse("http://api.duckduckgo.com/")
	parameters := url.Values{}
	parameters.Add("q", query)
	parameters.Add("format", "json")
	Url.RawQuery = parameters.Encode()
	var myClient = &http.Client{Timeout: 1 * time.Second}
	type duckduckgoResponse struct {
		RelatedTopics []struct {
			FirstURL string `json:"FirstURL"`
			Text     string `json:"Text"`
		} `json:"RelatedTopics"`
	}
	res, err := myClient.Get(Url.String())
	if err != nil {
		e <- &SearchError{source: "duckduckgo", message: "Request Timed out"}
		return
	}
	decoder := json.NewDecoder(res.Body)
	var response duckduckgoResponse
	for {
		if err := decoder.Decode(&response); err == io.EOF {
			break
		} else if err != nil {
			e <- &SearchError{source: "duckduckgo", message: "Error parsing JSON"}
			return
		}
	}
	var resultArray []SearchResult
	for _, element := range response.RelatedTopics {
		el := SearchResult{Text: element.Text, URL: element.FirstURL}
		resultArray = append(resultArray, el)
	}
	duckduckgoResults := Results{source: "duckduckgo", results: resultArray}
	r <- duckduckgoResults
}
