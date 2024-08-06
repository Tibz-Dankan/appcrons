package models

import (
	"encoding/json"
	"log"
	"time"
)

type PermittedRequestTime struct {
	ID string `json:"id"`
}

type PermittedApp struct {
	ID           string                 `json:"id"`
	RequestTimes []PermittedRequestTime `json:"requestTimes"`
}

type UserPermissions struct {
	UserID      string         `json:"userId"`
	Permissions []string       `json:"permissions"`
	Role        string         `json:"role"`
	Apps        []PermittedApp `json:"apps"`
}

type Permissions struct {
	User UserPermissions `json:"user"`
}

func (p *Permissions) Set(userId string) error {
	userPermissions := UserPermissions{UserID: userId}
	app := App{}
	user := User{}

	savedUser, err := user.FindOne(userId)
	if err != nil {
		log.Println("Error fetching user apps:", err)
		return err

	}

	if savedUser.ID == "" {
		log.Println("User does not exist!")
		return nil //TODO: To return custom error message
	}

	userPermissions.Role = savedUser.Role

	if userPermissions.Role == "client" {
		userPermissions.Permissions = []string{"read", "write", "update", "delete"}
	}

	if userPermissions.Role == "admin" {
		userPermissions.Permissions = []string{"read"}
	}

	userApps, err := app.FindByUserAndIncludeRequestTimes(userId)
	if err != nil {
		log.Println("Error fetching user apps:", err)
	}

	for _, app := range userApps {
		permittedApp := PermittedApp{}
		permittedApp.ID = app.ID

		if len(app.RequestTime) == 0 {
			userPermissions.Apps = append(userPermissions.Apps, permittedApp)
			continue
		}

		for _, requestTime := range app.RequestTime {
			permittedRequestTime := PermittedRequestTime{}
			permittedRequestTime.ID = requestTime.ID
			permittedApp.RequestTimes = append(permittedApp.RequestTimes, permittedRequestTime)
		}
		userPermissions.Apps = append(userPermissions.Apps, permittedApp)
	}

	if err := p.writeToCache(userPermissions); err != nil {
		return err
	}

	return nil
}

func (p *Permissions) Get(userId string) (UserPermissions, error) {
	permissions, err := p.readFromCache(userId)
	if err != nil {
		return permissions, err
	}

	return permissions, nil
}

func (p *Permissions) writeToCache(permissions UserPermissions) error {
	permissionJson, err := json.Marshal(&permissions)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return err
	}
	expiration := 9 * time.Hour

	if err = redisClient.Set(ctx, permissions.UserID, permissionJson, expiration).Err(); err != nil {
		log.Println("Error saving data to Redis:", err)
		return err
	}

	return nil
}

func (p *Permissions) readFromCache(userId string) (UserPermissions, error) {
	permissions := UserPermissions{}
	permissionString, err := redisClient.Get(ctx, userId).Result()
	if err != nil {
		log.Println("Error fetching data from Redis:", err)
		return permissions, nil
	}

	err = json.Unmarshal([]byte(permissionString), &permissions)
	if err != nil {
		log.Println("Error un-marshalling JSON:", err)
		return permissions, nil
	}

	return permissions, nil
}