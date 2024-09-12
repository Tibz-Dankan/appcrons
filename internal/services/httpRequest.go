package services

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
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
		log.Printf("Error creating request: %v", err)
		return response, err
	}

	req.Header.Set("Content-Type", "application/json")

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Minute * 4,
	}

	res, err := client.Do(req)
	duration := time.Since(startTime)
	requestTimeMS := int(duration.Milliseconds())

	if err != nil {
		log.Printf("Error making request: %v", err)
		if isTimeoutError(err) {
			retryBaseURL := os.Getenv("APPCRONS_RETRY_URL")
			retryURL := fmt.Sprintf("%s?url=%s", retryBaseURL, URL)

			if err := MakeHTTPRetryRequest(retryURL); err != nil {
				return response, err
			}
			response.StatusCode = 503
			response.RequestTimeMS = requestTimeMS
			response.StartedAt = startTime
			response.Message = "Request timeout"
			return response, nil
		}
		return response, err
	}

	response.StatusCode = res.StatusCode
	response.RequestTimeMS = requestTimeMS
	response.StartedAt = startTime

	return response, nil
}

func MakeHTTPRetryRequest(URL string) error {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 35,
	}

	res, err := client.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	log.Printf("Response Body: %s", string(body))

	return nil
}

func isTimeoutError(err error) bool {
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}
