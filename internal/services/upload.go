package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type Upload struct {
	StorageBucketBucket string `json:"storageBucket"`
	BaseURL             string `json:"baseURL"`
	DownloadURL         string `json:"url"`
	FilePath            string `json:"path"`
}

func (upload *Upload) initStorageBucket() (*storage.BucketHandle, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	currentDirPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	upload.BaseURL = "https://firebasestorage.googleapis.com/v0/b/"
	storageBucket := os.Getenv("STORAGE_BUCKET")
	upload.StorageBucketBucket = storageBucket

	configStorage := &firebase.Config{
		StorageBucket: storageBucket,
	}
	opt := option.WithCredentialsFile(currentDirPath + "/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), configStorage, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		return nil, err
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return nil, err
	}

	return bucket, nil

}

func (upload *Upload) Add(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {

	filePath := upload.FilePath
	if filePath == "" {
		return "", errors.New("no file path provided")
	}

	bucket, err := upload.initStorageBucket()
	if err != nil {
		return "", err
	}

	wc := bucket.Object(filePath).NewWriter(context.Background())
	_, err = io.Copy(wc, file)
	if err != nil {
		return "", err
	}

	err = wc.Close()
	if err != nil {
		return "", err
	}

	url, err := upload.getDownloadURL()
	if err != nil {
		return "", err
	}

	return url, nil
}

func (upload *Upload) Update(file multipart.File, fileHeader *multipart.FileHeader, savedFilePath string) (string, error) {

	filePath := upload.FilePath
	if filePath == "" {
		return "", errors.New("no file path provided")
	}

	bucket, err := upload.initStorageBucket()
	if err != nil {
		return "", err
	}

	wc := bucket.Object(filePath).NewWriter(context.Background())
	_, err = io.Copy(wc, file)
	if err != nil {
		return "", err
	}

	err = wc.Close()
	if err != nil {
		return "", err
	}

	url, err := upload.getDownloadURL()
	if err != nil {
		return "", err
	}

	if url != "" {
		if err := upload.Delete(savedFilePath); err != nil {
			return "", err
		}

		fmt.Println("file deleted from storage using path ==>", savedFilePath)
	}

	return url, nil
}

func (upload *Upload) Delete(filePath string) error {

	if filePath == "" {
		return errors.New("no file path provided")
	}

	bucket, err := upload.initStorageBucket()
	if err != nil {
		return err
	}
	obj := bucket.Object(filePath)

	if err := obj.Delete(context.Background()); err != nil {
		return err
	}

	return nil
}

func (upload *Upload) transformFilePath() (string, error) {
	path := upload.FilePath

	if path == "" {
		return "", errors.New("no file path provided")
	}

	path = strings.ReplaceAll(path, "/", "%2F")
	path = strings.ReplaceAll(path, " ", "%20")
	path = strings.ReplaceAll(path, "?", "%3F")
	path = strings.ReplaceAll(path, "&", "%26")
	path = strings.ReplaceAll(path, "=", "%3D")
	path = strings.ReplaceAll(path, ":", "%3A")
	path = strings.ReplaceAll(path, ",", "%2C")

	return path, nil
}

func (upload *Upload) getDownloadURL() (string, error) {
	start := time.Now()

	transformedFilePath, err := upload.transformFilePath()
	if err != nil {
		return "", err
	}

	FIREBASE_STORAGE_URL := upload.BaseURL + upload.StorageBucketBucket + "/o/" + transformedFilePath

	req, err := http.NewRequest(http.MethodGet, FIREBASE_STORAGE_URL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	type Response struct {
		Name               string `json:"name"`
		Bucket             string `json:"bucket"`
		Generation         string `json:"generation"`
		Metageneration     string `json:"metageneration"`
		ContentType        string `json:"contentType"`
		TimeCreated        string `json:"timeCreated"`
		Updated            string `json:"updated"`
		StorageClass       string `json:"storageClass"`
		Size               string `json:"size"`
		Md5Hash            string `json:"md5Hash"`
		ContentEncoding    string `json:"contentEncoding"`
		ContentDisposition string `json:"contentDisposition"`
		Crc32c             string `json:"crc32c"`
		Etag               string `json:"etag"`
		DownloadTokens     string `json:"downloadTokens"`
	}

	if res.StatusCode != http.StatusOK {
		_, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		return "", errors.New("request to firebase storage failed")
	}

	fmt.Printf("firebase-storage-service: status code: %d\n", res.StatusCode)
	rBody, _ := io.ReadAll(res.Body)

	response := Response{}
	json.NewDecoder(strings.NewReader(string(rBody))).Decode(&response)

	downloadURL := FIREBASE_STORAGE_URL + "?alt=media&token=" + response.DownloadTokens
	upload.DownloadURL = downloadURL

	fmt.Printf(
		"Firebase storage request duration : %s\n",
		time.Since(start),
	)

	return downloadURL, nil
}
