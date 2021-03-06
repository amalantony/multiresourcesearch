package search

import (
	"strconv"

	"time"

	"github.com/ChimeraCoder/anaconda"
)

const (
	consumerKey       = "5qn8ZqG3tyned7TiN0o1xU5nN"
	consumerSecret    = "h99EAfvLTXkpiJZgaAtmNOsVIU0ipWGAQoS6yUWWARI7ST1Sw3"
	accessToken       = "19109487-WUKiMTWJqnFDDN8RLAXFU0rkMXNQSvgm0prziL0XS"
	accessTokenSecret = "d96BCYhsQcHfPtnraj0YyJSGi1dGTyXVVHtCNfk7m834O"
)

func getTwitterResults(query string, r chan Results, e chan error) {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	ch := make(chan Results)
	ech := make(chan error)
	after := time.After(1 * time.Second)
	go func(ch chan Results, ech chan error) {
		results, err := api.GetSearch(query, nil)
		if err != nil {
			ech <- &SearchError{source: "twitter", message: "Error making Request"}
			return
		}
		var resultArray []SearchResult
		for _, tweet := range results.Statuses {
			url := "https://twitter.com/statuses/" + strconv.FormatInt(tweet.Id, 10)
			el := SearchResult{Text: tweet.Text, URL: url}
			resultArray = append(resultArray, el)
		}
		twitterResults := Results{source: "twitter", results: resultArray}
		ch <- twitterResults
	}(ch, ech)
	select {
	case <-after: // request timed out
		e <- &SearchError{source: "twitter", message: "Request Timed out"}
	case twitterResults := <-ch:
		r <- twitterResults
	case twitterError := <-ech:
		e <- twitterError
	}
}
