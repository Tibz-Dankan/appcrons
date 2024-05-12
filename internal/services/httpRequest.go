package services

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
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

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	// Create a cancellable context and wire it up to signals from Ctrl-C.
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-signals
		log.Println("Request cancelled with Ctrl-C")
		cancel()
	}()

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return response, err
	}

	req.Header.Set("Content-Type", "application/json")
	defaultClient := http.DefaultClient
	defaultClient.Timeout = time.Second * 30
	res, err := defaultClient.Do(req.WithContext(ctx))

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

	log.Printf("Request statusCode: %d Duration: %d URL: %s\n", res.StatusCode, requestTimeMS, URL)

	return response, nil
}

func isTimeoutError(err error) bool {
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}
