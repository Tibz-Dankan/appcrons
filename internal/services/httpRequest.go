package services

import (
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
	defaultClient := http.DefaultClient
	defaultClient.Timeout = time.Second * 30
	res, err := defaultClient.Do(req)

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
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}
