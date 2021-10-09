package cmd

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"log"
	"net/http"
)

func HandleBing(w http.ResponseWriter, r *http.Request) {
	decoder := schema.NewDecoder()
	var sQuery SearchQuery
	err := decoder.Decode(&sQuery, r.URL.Query())
	if err != nil {
		http.Error(w, "Failed to parse query", 400)
		log.Println("Failed to parse query")
		return
	}

	result, err := Bing(sQuery)

	if err != nil {
		http.Error(w, "Failed to query api", 503)
		log.Println("Failed to query api ", err.Error())
		return
	}

	payload, err := json.Marshal(result)
	if err != nil {
		log.Println("Failed to marshal result")
	}

	_, err = w.Write(payload)
	if err != nil {
		log.Println("Failed to write using ResponseWriter")
	}
}
