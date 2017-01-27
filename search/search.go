package search

import "fmt"

// SearchResult ... individual search result
type SearchResult struct {
	URL, Text string
}

// Results ... collection of searchResults from a source
type Results struct {
	source  string
	results []SearchResult
}

// SearchError ... error implementation
type SearchError struct {
	source, message string
}

// AggregatedResults ... store the aggreated search results
type AggregatedResults struct {
	Query  string
	Data   map[string][]SearchResult
	Errors []string
}

func (e *SearchError) Error() string {
	return fmt.Sprintf("%s - %s", e.source, e.message)
}

// Search ... perform the multi resource parallel search
func Search(query string) AggregatedResults {
	r := make(chan Results)
	e := make(chan error)

	go getDuckDuckGoResults(query, r, e)
	go getGoogleResults(query, r, e)
	go getTwitterResults(query, r, e)

	aggregatedResults := AggregatedResults{Query: query, Data: make(map[string][]SearchResult)}

	cnt := 0

	for {
		select {
		case searchResults := <-r:
			aggregatedResults.Data[searchResults.source] = searchResults.results
			cnt++
		case errors := <-e:
			aggregatedResults.Errors = append(aggregatedResults.Errors, errors.Error())
			cnt++
		}
		if cnt == 3 { // break out after getting values/errors from each of the 3 go-routines
			return aggregatedResults
		}
	}
}
