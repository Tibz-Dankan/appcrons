package models

import (
	"encoding/json"
	"log"
	"time"
)

type AppCache struct{}

func (uc *AppCache) Write(app App) error {
	// Convert struct to JSON.
	appData, err := json.Marshal(&app)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return err
	}

	expiration := 3 * time.Hour

	if err = redisClient.Set(ctx, app.ID, appData, expiration).Err(); err != nil {
		log.Println("Error saving data to Redis:", err)
		return err
	}

	return nil
}

func (ac *AppCache) WriteByUser(userId string, apps []App) error {
	// Convert struct to JSON.
	appsData, err := json.Marshal(&apps)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return err
	}

	expiration := 3 * time.Hour
	var key = "user-apps:" + userId

	if err = redisClient.Set(ctx, key, appsData, expiration).Err(); err != nil {
		log.Println("Error saving data to Redis:", err)
		return err
	}

	return nil
}

func (ac *AppCache) Read(key string) (App, error) {

	app := App{}
	savedAppData, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		log.Println("Error fetching data from Redis:", err)
		return app, nil
	}

	// Convert string into JSON.
	err = json.Unmarshal([]byte(savedAppData), &app)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return app, nil
	}

	return app, nil
}

func (ac *AppCache) ReadByUser(userId string) ([]App, error) {

	apps := []App{}
	var key = "user-apps:" + userId

	savedAppsData, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		log.Println("Error fetching data from Redis:", err)
		return apps, nil
	}

	// Convert string into JSON.
	err = json.Unmarshal([]byte(savedAppsData), &apps)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return apps, nil
	}

	return apps, nil
}

func (ac *AppCache) WriteAll(apps []App) error {
	// Convert struct to JSON.
	appsData, err := json.Marshal(&apps)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return err
	}

	expiration := 3 * time.Hour
	var key = "apps"

	if err = redisClient.Set(ctx, key, appsData, expiration).Err(); err != nil {
		log.Println("Error saving data to Redis:", err)
		return err
	}

	return nil
}

func (ac *AppCache) ReadAll() ([]App, error) {
	apps := []App{}
	var key = "apps"

	savedAppsData, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		log.Println("Error fetching data from Redis:", err)
		return apps, nil
	}

	// Convert string into JSON.
	err = json.Unmarshal([]byte(savedAppsData), &apps)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return apps, nil
	}

	return apps, nil
}

func (ac *AppCache) WriteOneToAll(app App) error {
	var apps []App
	var key = "apps"
	// get all existing apps in the cache
	savedApps, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		log.Println("Error fetching data from Redis:", err)
		return err
	}

	// Convert struct into JSON.
	err = json.Unmarshal([]byte(savedApps), &apps)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return err
	}

	// update the app info in the app array
	for i, a := range apps {
		if a.ID == app.ID {
			apps[i] = app
		}
	}

	// Convert struct to JSON.
	appsData, err := json.Marshal(&apps)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return err
	}

	expiration := 3 * time.Hour
	// save the updated apps
	if err = redisClient.Set(ctx, key, appsData, expiration).Err(); err != nil {
		log.Println("Error saving data to Redis:", err)
		return err
	}

	return nil
}

func (ac *AppCache) WriteOneToUser(app App) error {
	apps := []App{}
	var key = "user-apps:" + app.UserID
	// get all existing apps in the cache
	savedApps, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		log.Println("Error fetching data from Redis:", err)
		return err
	}

	// Convert struct into JSON.
	err = json.Unmarshal([]byte(savedApps), &apps)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return err
	}

	// update the app info in the app array
	for i, a := range apps {
		if a.ID == app.ID {
			apps[i] = app
		}
	}

	// Convert struct to JSON.
	appsData, err := json.Marshal(&apps)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return err
	}

	expiration := 3 * time.Hour
	// save the updated apps
	if err = redisClient.Set(ctx, key, appsData, expiration).Err(); err != nil {
		log.Println("Error saving data to Redis:", err)
		return err
	}

	return nil
}

func (ac *AppCache) Delete(appID string) error {
	// Delete data from Redis
	err := redisClient.Del(ctx, appID).Err()
	if err != nil {
		log.Println("Error deleting data from Redis:", err)
		return err
	}

	return nil
}
