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
	userApps := []App{}
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

	for _, app := range apps {
		if app.UserID == userId {
			userApps = append(userApps, app)
		}
	}

	return userApps, nil
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

func (ac *AppCache) hasApp(apps []App, id string) bool {
	for _, app := range apps {
		if app.ID == id {
			return true
		}
	}
	return false
}

func (ac *AppCache) WriteOneToAll(app App) error {
	apps, err := ac.ReadAll()
	if err != nil {
		return err
	}

	// Update Existing app and cache
	if appFound := ac.hasApp(apps, app.ID); appFound {
		for i, a := range apps {
			if a.ID == app.ID {
				apps[i] = app
			}
		}
		if err := ac.WriteAll(apps); err != nil {
			return err
		}
		return nil
	}

	// Add new app and update cache
	apps = append(apps, app)
	if err := ac.WriteAll(apps); err != nil {
		return err
	}

	return nil
}

func (ac *AppCache) WriteOneToUser(app App) error {
	apps, err := ac.ReadByUser(app.UserID)
	if err != nil {
		return err
	}

	// Update Existing app and cache
	if appFound := ac.hasApp(apps, app.ID); appFound {
		for i, a := range apps {
			if a.ID == app.ID {
				apps[i] = app
			}
		}
		if err := ac.WriteAll(apps); err != nil {
			return err
		}
		return nil
	}

	// Add new app and update cache
	apps = append(apps, app)
	if err := ac.WriteByUser(app.UserID, apps); err != nil {
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
