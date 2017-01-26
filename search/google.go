package search

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func getGoogleResults(query string, r chan Results, e chan error) {
	url := "https://www.googleapis.com/customsearch/v1?key=AIzaSyD5bMRAkJmNSgE-TtsW21uTZMDMMNPE0rM&cx=015968340880156748630%3Aui5gmgfozyu&q=" + query
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
	dec := json.NewDecoder(res.Body)
	var response googleResponse
	for {
		if err := dec.Decode(&response); err == io.EOF {
			break
		} else if err != nil {
			e <- &SearchError{source: "google", message: "Error parsing JSON"}
			return
		}
	}
	resultArray := make([]searchResult, len(response.Items))
	for _, element := range response.Items {
		el := searchResult{text: element.Snippet, url: element.Link}
		if len(el.url) > 0 {
			resultArray = append(resultArray, el)
		}
	}
	googleResults := Results{source: "google", results: resultArray}
	r <- googleResults
}
