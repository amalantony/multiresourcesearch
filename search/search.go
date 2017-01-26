package search

import (
	"fmt"
	"time"
)

type searchResult struct { // individual search result
	url, text string
}

// Results ... Results is an array of the result type
type Results struct {
	source  string
	results []searchResult
}

// Search ... perform the multi resource parallel search
func Search(query string) {
	r := make(chan Results)
	e := make(chan error)

	go getDuckDuckGoResults(query, r, e)

	cnt := 0
	for {
		select {
		case searchResults := <-r:
			fmt.Println(searchResults)
			cnt++
		case errors := <-e:
			fmt.Println("Request timed out", errors)
			cnt++
		}
		if cnt == 1 {
			break
		}
	}

	time.Sleep(time.Second) // for debugging print
}
