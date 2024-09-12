package services

import (
	"log"
	"net"
	"net/http"
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
		if isTimeoutError(err) {
			response.StatusCode = 503
			response.RequestTimeMS = requestTimeMS
			response.StartedAt = startTime
			return response, nil
		}
		return response, err
	}

	response.StatusCode = res.StatusCode
	response.RequestTimeMS = requestTimeMS
	response.StartedAt = startTime

	return response, nil
}

func isTimeoutError(err error) bool {
	log.Println("Timeout error: ", err)
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}
