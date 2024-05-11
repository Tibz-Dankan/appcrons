package services

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Message       string    `json:"message"`
	StatusCode    int       `json:"statusCode"`
	RequestTimeMS int       `json:"requestTimeMS"`
	StartedAt     time.Time `json:"startedAt"`
}

func MakeHTTPRequest(URL string) (Response, error) {
	response := Response{}
	startTime := time.Now()

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return response, err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return response, err
	}

	duration := time.Since(startTime)
	requestTimeMS := int(duration.Milliseconds())

	type Body struct {
		Message string `json:"message"`
	}

	body := Body{}
	resBody, _ := io.ReadAll(res.Body)
	json.NewDecoder(strings.NewReader(string(resBody))).Decode(&body)

	response.StatusCode = res.StatusCode
	response.Message = body.Message
	response.RequestTimeMS = requestTimeMS
	response.StartedAt = startTime

	log.Printf("Request statusCode: %d Duration: %d URL: %s\n", res.StatusCode, requestTimeMS, URL)
	log.Printf("Response body: %s\n", resBody)

	return response, nil
}
