package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type RequestCache struct{}

func (rc *RequestCache) Write(request Request) error {
	// Convert struct to JSON.
	requestData, err := json.Marshal(&request)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	expiration := 5 * time.Minute

	if err = redisClient.Set(ctx, request.ID, requestData, expiration).Err(); err != nil {
		fmt.Println("Error saving data to Redis:", err)
		return err
	}

	return nil
}

func (rc *RequestCache) WriteByApp(appId string, requests []Request) error {
	// Convert struct to JSON.
	appRequestData, err := json.Marshal(&requests)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	expiration := 5 * time.Minute
	var key = "app-requests:" + appId

	if err = redisClient.Set(ctx, key, appRequestData, expiration).Err(); err != nil {
		fmt.Println("Error saving data to Redis:", err)
		return err
	}

	return nil
}

func (rc *RequestCache) Read(key string) (Request, error) {

	request := Request{}
	savedAppData, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("Error fetching data from Redis:", err)
		return request, nil
	}

	// Convert string into JSON.
	err = json.Unmarshal([]byte(savedAppData), &request)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return request, nil
	}

	return request, nil
}

func (rc *RequestCache) ReadByApp(appId string) ([]Request, error) {

	apps := []Request{}
	var key = "app-requests:" + appId

	savedAppsData, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("Error fetching data from Redis:", err)
		return apps, nil
	}

	// Convert string into JSON.
	err = json.Unmarshal([]byte(savedAppsData), &apps)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return apps, nil
	}

	return apps, nil
}

func (rc *RequestCache) Delete(requestId string) error {
	// Delete data from Redis
	err := redisClient.Del(ctx, requestId).Err()
	if err != nil {
		fmt.Println("Error deleting data from Redis:", err)
		return err
	}

	return nil
}
