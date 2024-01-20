package main

import (
	// "context"
	// "fmt"
	// "log"
	// "valentin-lvov/1x-parser/scrapper"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
)

type TrackRequest struct {
	URL      string `json:"url"`
	Duration int    `json:"duration"`
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func trackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method allowed on this endpoint", http.StatusMethodNotAllowed)
		return
	}
	var request TrackRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token := GenerateSecureToken(20)
	// TODO: add rabbitMQ integration
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	/*endpoint looks like this: http://example.com/api/results?token=12345*/
	if r.Method != "GET" {
		http.Error(w, "Only GET requests on this endpoint", http.StatusMethodNotAllowed)
		return
	}

	var token string
	var results map[string]string
	token = r.URL.Query().Get("token")

	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	// TODO: Retrieve tracking results from the database or cache
	// results, err := getTrackingResults(token)
	// if err != nil {
	//     http.Error(w, err.Error(), http.StatusInternalServerError)
	//     return
	// }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]map[string]string{"data": results}) // Replace "results" with actual data

}
func main() {
	/*
		// EXAMPLE OF JUST SCRAPPING PARTUCULAR PAGE
		url := "https://1xbet.com/en/live/football/21317-africa-cup-of-nations/504058217-cape-verde-mozambique"

		var result map[string]string
		var ctx *context.Context
		var err error

		ctx, cancel, err := scrapper.MakeConnectionAndLoad(url)
		defer cancel()

		if err != nil {
			log.Fatal("Error creating ChromeDP context:", err)
			return
		}
		result, err = scrapper.GetContent(ctx, "div.bet-inner")
		if err != nil {
			log.Fatal("Error getting the content:", err)
			return
		}

		fmt.Println(len(result))
	*/

	http.HandleFunc("/api/track", trackHandler)
	http.HandleFunc("/api/results", resultsHandler)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
