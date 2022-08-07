package cmd

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const v7Endpoint string = "https://api.bing.microsoft.com/v7.0/search"

var apiKey = os.Getenv("ocp_apim_subscription_key")

func Bing(search SearchQuery, queryStringAnalyzerChannel chan<- string) (*BingAnswer, error) {
	// Declare a new GET request.
	req, err := http.NewRequest("GET", v7Endpoint, nil)
	if err != nil {
		log.Println("Failed to create HTTP request.")
		return nil, err
	}

	// Add the payload to the request.
	param := req.URL.Query()
	param.Add("q", search.Query)
	req.URL.RawQuery = param.Encode()

	// Insert the request header.
	req.Header.Add("Ocp-Apim-Subscription-Key", apiKey)

	// Instantiate a client.
	client := &http.Client{}

	// Send the request to Bing.
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send request to ", v7Endpoint)
		return nil, err
	}

	// send query text to channel for processing
	queryStringAnalyzerChannel <- search.Query

	// defer close the response when code is out of scope.
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Failed to close io reader")
		}
	}(resp.Body)

	// read body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to parse result body")
		return nil, err
	}

	// Create a new BingAnswer
	ans := BingAnswer{}
	err = json.Unmarshal(body, &ans)
	if err != nil {
		log.Println("Failed to unmarshal to BingResult struct")
		return nil, err
	}

	return &ans, nil
}
