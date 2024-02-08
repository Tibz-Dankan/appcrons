package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Message       string `json:"message"`
	StatusCode    int    `json:"statusCode"`
	RequestTimeMS int    `json:"requestTimeMS"`
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

	fmt.Printf("Request status code: %d\n", res.StatusCode)
	fmt.Printf("Response body: %s\n", resBody)

	return response, nil
}
