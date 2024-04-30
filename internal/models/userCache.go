package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// ======CONSIDER CHANGING WAY OF IMPLEMENTING REDIS JSON=====
// ======CONSIDER CHANGING WAY OF IMPLEMENTING REDIS JSON=====

type UserCache struct{}

func (uc *UserCache) Write(user User) error {
	fmt.Println("writing to cache....")

	// Convert struct to JSON.
	userData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	expiration := 3 * time.Hour / time.Second

	fmt.Println("User to be saved in the cache:", user)

	err = redisClient.Set(ctx, user.ID, userData, expiration).Err()
	if err != nil {
		fmt.Println("Error saving data to Redis:", err)
		return err
	}

	return nil
}

func (uc *UserCache) Read(userId string) (User, error) {
	fmt.Println("Reading from cache....")

	user := User{}
	savedUserData, err := redisClient.Get(ctx, userId).Result()
	if err != nil {
		fmt.Println("Error fetching data from Redis:", err)
		return user, nil
	}

	// Convert JSON into struct.
	err = json.Unmarshal([]byte(savedUserData), &user)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return user, nil
	}

	return user, nil
}

func (uc *UserCache) Delete(userID string) error {
	// Delete data from Redis
	err := redisClient.Del(ctx, userID).Err()
	if err != nil {
		fmt.Println("Error deleting data from Redis:", err)
		return err
	}

	return nil
}
