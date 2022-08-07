package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/schema"
)

func HandleBing(analyzer chan<- string, w http.ResponseWriter, r *http.Request) {

	// decode into a SearchQuery
	decoder := schema.NewDecoder()
	var sQuery SearchQuery
	err := decoder.Decode(&sQuery, r.URL.Query())
	if err != nil {
		http.Error(w, "Failed to parse query", 400)
		log.Println("Failed to parse query")
		return
	}

	// get the BingAnswer
	result, err := Bing(sQuery, analyzer)

	// api result
	if err != nil {
		http.Error(w, "Failed to query api", http.StatusServiceUnavailable)
		log.Println("Failed to query api ", err.Error())
		return
	}

	// deserialization
	payload, err := json.Marshal(result)
	if err != nil {
		log.Println("Failed to marshal result")
	}

	_, err = w.Write(payload)
	if err != nil {
		log.Println("Failed to write using ResponseWriter")
	}
}
