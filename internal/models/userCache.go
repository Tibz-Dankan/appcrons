package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type UserCache struct{}

func (uc *UserCache) Write(user User) error {
	// Convert struct to JSON.
	userData, err := json.Marshal(&user)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	expiration := 3 * time.Hour

	if err = redisClient.Set(ctx, user.ID, userData, expiration).Err(); err != nil {
		fmt.Println("Error saving data to Redis:", err)
		return err
	}
	if err = redisClient.Set(ctx, user.Email, userData, expiration).Err(); err != nil {
		fmt.Println("Error saving data to Redis:", err)
		return err
	}

	return nil
}

func (uc *UserCache) Read(key string) (User, error) {

	user := User{}
	savedUserData, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("Error fetching data from Redis:", err)
		return user, nil
	}

	// Convert string into JSON.
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
