package services

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type FirebaseAdmin struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

type FirebaseManager struct {
	fileName string
	data     FirebaseAdmin
}

func NewFirebaseManager(fileName string) *FirebaseManager {

	env := os.Getenv("GO_ENV")
	if env == "development" || env == "testing" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
		log.Println("Loaded env var file")
	}

	return &FirebaseManager{
		fileName: fileName,
		data: FirebaseAdmin{
			Type:                    os.Getenv("FIREBASE_TYPE"),
			ProjectID:               os.Getenv("FIREBASE_PROJECT_ID"),
			PrivateKeyID:            os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
			PrivateKey:              formatPrivateKey(os.Getenv("FIREBASE_PRIVATE_KEY")),
			ClientEmail:             os.Getenv("FIREBASE_CLIENT_EMAIL"),
			ClientID:                os.Getenv("FIREBASE_CLIENT_ID"),
			AuthURI:                 os.Getenv("FIREBASE_AUTH_URI"),
			TokenURI:                os.Getenv("FIREBASE_TOKEN_URI"),
			AuthProviderX509CertURL: os.Getenv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL"),
			ClientX509CertURL:       os.Getenv("FIREBASE_CLIENT_X509_CERT_URL"),
			UniverseDomain:          os.Getenv("FIREBASE_UNIVERSE_DOMAIN"),
		},
	}
}

func (fm *FirebaseManager) CreateFile() error {
	startTime := time.Now()
	file, err := os.Create(fm.fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(fm.data, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}
	log.Println("createFileDuration:", time.Since(startTime))

	log.Println(fm.fileName, "file created successfully!")
	return nil
}

func (fm *FirebaseManager) DeleteFile() error {
	startTime := time.Now()
	err := os.Remove(fm.fileName)
	if err != nil {
		return err
	}

	log.Println("deleteFileDuration:", time.Since(startTime))

	log.Println(fm.fileName, "file deleted successfully!")
	return nil
}

func formatPrivateKey(privateKey string) string {
	privateKey = strings.TrimSpace(privateKey)

	privateKey = strings.ReplaceAll(privateKey, "-----BEGIN PRIVATE KEY-----", "")
	privateKey = strings.ReplaceAll(privateKey, "-----END PRIVATE KEY-----", "")

	privateKey = strings.ReplaceAll(privateKey, "\\n", "")

	var lines []string
	for i := 0; i < len(privateKey); i += 64 {
		end := i + 64
		if end > len(privateKey) {
			end = len(privateKey)
		}
		lines = append(lines, privateKey[i:end])
	}

	processedKey := "-----BEGIN PRIVATE KEY-----\n" +
		strings.Join(lines, "\n") +
		"\n-----END PRIVATE KEY-----\n"

	return processedKey
}
