package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

func main() {}

// Handler is the fission handler
func Handler(w http.ResponseWriter, r *http.Request) {

	log.SetOutput(os.Stdout)

	redisClient, err := client("binging-parrot-redis:6379", "wQC4kdylua")
	if err != nil {
		http.Error(w, fmt.Sprintf("could not connect to REDIS: %v", err), http.StatusInternalServerError)
		return
	}

	httpBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not read request body: %v", err), http.StatusInternalServerError)
		return
	}

	payload := make(map[string]interface{}, 0)
	payload["timestamp_ns"] = time.Now().UnixNano()
	payload["request_body"] = httpBody

	jsonBytes, err := json.Marshal(&payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal data structure: %v", err), http.StatusInternalServerError)
		return
	}

	if err := saveAnalytics(redisClient, jsonBytes); err != nil {
		http.Error(w, fmt.Sprintf("failed to store request: %v", err), http.StatusInternalServerError)
		return
	}

	return
}

func client(address, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       0,        // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, err
}

func saveAnalytics(client *redis.Client, payload []byte) error {
	id := uuid.New()
	status := client.Set(id.String(), payload, 24*time.Hour)
	return status.Err()
}
