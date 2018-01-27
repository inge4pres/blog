package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

func main() {

	log.SetOutput(os.Stdout)

	redisClient, err := client(os.Getenv("REDIS_SERVER_ADDR"), os.Getenv("REDIS_SERVER_PWD"))
	if err != nil {
		log.Fatalf("could not connect to REDIS: %v", err)
	}

	httpBody, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("could not read from Stdin: %v", err)
	}

	payload := make(map[string]interface{}, 0)
	payload["timestamp_ns"] = time.Now().UnixNano()
	payload["request_body"] = httpBody

	jsonBytes, err := json.Marshal(&payload)
	if err != nil {
		log.Printf("failed to marshal data structure: %v", err)
		return
	}

	if err := saveAnalytics(redisClient, jsonBytes); err != nil {
		log.Printf("failed to store request: %v", err)
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
