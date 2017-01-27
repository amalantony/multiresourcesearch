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
	after := time.After(1 * time.Second)
	var resultArray []searchResult
	go func(ch chan Results) {
		results, _ := api.GetSearch(query, nil)
		resultArray = make([]searchResult, len(results.Statuses))
		for _, tweet := range results.Statuses {
			url := "https://twitter.com/statuses/" + strconv.FormatInt(tweet.Id, 10)
			el := searchResult{text: tweet.Text, url: url}
			resultArray = append(resultArray, el)
		}
		twitterResults := Results{source: "twitter", results: resultArray}
		ch <- twitterResults
	}(ch)
	select {
	case <-after: //timedout
		e <- &searchError{source: "twitter", message: "Request timed out!"}
	case twitterResults := <-ch: //
		r <- twitterResults
	}
}
