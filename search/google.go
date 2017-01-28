package search

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	apiKey               = "AIzaSyD5bMRAkJmNSgE-TtsW21uTZMDMMNPE0rM"
	customSearchEngineID = "015968340880156748630%3Aui5gmgfozyu"
)

func getGoogleResults(query string, r chan Results, e chan error) {
	Url, _ := url.Parse("https://www.googleapis.com/customsearch/v1")
	parameters := url.Values{}
	parameters.Add("key", apiKey)
	parameters.Add("q", query)
	Url.RawQuery = parameters.Encode()
	var myClient = &http.Client{Timeout: 1 * time.Second}
	type googleResponse struct {
		Items []struct {
			Link    string `json:"link"`
			Snippet string `json:"snippet"`
		} `json:"items"`
	}
	res, err := myClient.Get(Url.String() + "&cx=" + customSearchEngineID) // due to special characeters in the customSearchEngineID
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
