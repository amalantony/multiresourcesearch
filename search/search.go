package search

import (
	"fmt"
	"time"
)

type searchResult struct { // individual search result
	url, text string
}

// Results ... Results is a collection of searchResults from a particular source
type Results struct {
	source  string
	results []searchResult
}

// SearchError ... Error type for searches
type searchError struct {
	source, message string
}

func (e *searchError) Error() string {
	return fmt.Sprintf("%s - %s", e.source, e.message)
}

// Search ... perform the multi resource parallel search
func Search(query string) {
	r := make(chan Results)
	e := make(chan error)

	go getDuckDuckGoResults(query, r, e)
	go getGoogleResults(query, r, e)
	go getTwitterResults(query, r, e)

	cnt := 0
	for {
		select {
		case searchResults := <-r:
			// look at searchResults.source and construct response
			fmt.Println(searchResults)
			cnt++
		case errors := <-e:
			// look at error type & construct response
			fmt.Println(errors)
			cnt++
		}
		if cnt == 3 { // break out after getting values/errors from each of the 3 go-routines
			break
		}
	}

	time.Sleep(time.Second) // for debugging print, remove in production
}
