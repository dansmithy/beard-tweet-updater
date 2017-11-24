package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
)

const defaultRootDir = "/appdata"

const twitterAppURL = "http://www.beards-and-other-facial-hair-styles.co.uk"

func latestTweet() (*twittergo.Tweet, error) {
	var (
		err     error
		client  *twittergo.Client
		req     *http.Request
		resp    *twittergo.APIResponse
		results *twittergo.SearchResults
	)

	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_KEY")

	if consumerKey == "" || consumerSecret == "" {
		return nil, errors.New("TWITTER_CONSUMER_KEY and/or TWITTER_CONSUMER_SECRET not set")
	}

	config := &oauth1a.ClientConfig{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
	}
	client = twittergo.NewClient(config, nil)

	query := url.Values{}
	query.Set("q", "beard")
	query.Set("count", "1")
	query.Set("result_type", "recent")
	url := fmt.Sprintf("/1.1/search/tweets.json?%v", query.Encode())

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("could not parse request")
	}

	resp, err = client.SendRequest(req)
	if err != nil {
		return nil, errors.New("could not send request")
	}

	results = &twittergo.SearchResults{}
	err = resp.Parse(results)
	if err != nil {
		return nil, errors.New("problem parsing response")
	}

	if len(results.Statuses()) == 0 {
		return nil, errors.New("no matching tweets")
	}
	return &results.Statuses()[0], nil
}

func recordExecutionTimestamp(rootDir string) {
	ioutil.WriteFile(rootDir+"/last_run_date", []byte(time.Now().Format("2006-01-02-150405")+"\n"), 0644)
}

func recordTweet(rootDir string, tweetText string) {
	ioutil.WriteFile(rootDir+"/last_tweet", []byte(tweetText+"\n"), 0644)
}

func main() {

	var rootDir = defaultRootDir
	if len(os.Args) > 1 {
		rootDir = os.Args[1]
	}

	recordExecutionTimestamp(rootDir)
	tweet, err := latestTweet()

	if tweet != nil {
		recordTweet(rootDir, tweet.Text())
	}

	if err != nil {
		fmt.Println("Had an error:", err)
	}

}
