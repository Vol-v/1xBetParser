package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	"fmt"
	"log"
	"net/http"
	"valentin-lvov/1x-parser/cache"
	"valentin-lvov/1x-parser/queue"

	// "valentin-lvov/1x-parser/scrapper"

	"github.com/redis/go-redis/v9"
)

type TrackRequest struct {
	URL      string `json:"url"`
	Duration int    `json:"duration"`
}

var rdb *redis.Client

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

	err = queue.PublishTrackingTask(request.URL, request.Duration)

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]string{"token": token})
	// token := GenerateSecureToken(20)

	w.WriteHeader(http.StatusOK)
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	/*endpoint looks like this: http://example.com/api/results?url=12345*/
	if r.Method != "GET" {
		http.Error(w, "Only GET requests on this endpoint", http.StatusMethodNotAllowed)
		return
	}

	var url string
	var results map[string]string
	url = r.URL.Query().Get("url")

	if url == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	// TODO: Retrieve tracking results from the database or cache
	results, err := cache.RetrieveFromRedis(rdb, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]map[string]string{"data": results}) // Replace "results" with actual data

}
func main() {

	// EXAMPLE OF JUST SCRAPPING PARTUCULAR PAGE
	// url := "https://1xbet.com/en/live/football/96463-germany-bundesliga/504601794-borussia-monchengladbach-augsburg"

	// var result map[string]string
	// var ctx *context.Context
	// var err error

	// ctx, cancel, err := scrapper.MakeConnectionAndLoad(url)
	// defer cancel()

	// if err != nil {
	// 	log.Fatal("Error creating ChromeDP context:", err)
	// 	return
	// }
	// result, err = scrapper.GetContentFromSelector(ctx, "div.bet-inner")
	// if err != nil {
	// 	log.Fatal("Error getting the content:", err)
	// 	return
	// }

	// fmt.Println(len(result))

	rdb = cache.NewRedisClient()
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println(pong) // Output: PONG (if successful)

	go queue.StartConsumer(rdb)

	http.HandleFunc("/api/track", trackHandler)
	http.HandleFunc("/api/results", resultsHandler)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
