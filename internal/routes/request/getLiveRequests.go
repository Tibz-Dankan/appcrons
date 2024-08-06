package request

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/event"
	"github.com/Tibz-Dankan/keep-active/internal/middlewares"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

type AppRequestProgress = services.AppRequestProgress

func sendMessage(message, userId string, clientManager *services.ClientManager) error {
	w, ok := clientManager.GetClient(userId)
	if !ok {
		log.Println("Client not found")
		return nil
	}

	f, ok := w.(http.Flusher)
	if !ok {
		log.Println("Response writer does not implement http.Flusher")
		return nil
	}

	data, _ := json.Marshal(map[string]string{
		"message": message,
		"userId":  userId,
	})

	_, err := w.Write([]byte("data: " + string(data) + "\n\n"))
	if err != nil {
		log.Println("Error writing to response writer:", err)
		return err
	}
	f.Flush()
	return nil
}

func sendAppToClient(appRequestProgress AppRequestProgress, clientManager *services.ClientManager) error {
	client, ok := clientManager.GetClient(appRequestProgress.UserID)
	if !ok {
		log.Println("Client not found")
		return nil
	}

	appJson, err := json.Marshal(&appRequestProgress)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return err
	}

	f, ok := client.(http.Flusher)
	if !ok {
		log.Println("Client does not implement http.Flusher")
		return err
	}
	_, err = client.Write([]byte("data: " + string(appJson) + "\n\n"))
	if err != nil {
		log.Println("Error sending event:", err)
		return err
	}
	f.Flush()

	return nil
}

func getLiveRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	userId, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		services.AppError("UserID not found in context", 500, w)
		return
	}
	log.Println("Client connected:", userId)

	clientManager := services.NewClientManager()
	clientManager.AddClient(userId, w)

	if err := sendMessage("warmup", userId, clientManager); err != nil {
		return
	}

	type AppRequestProgress = services.AppRequestProgress
	appCh := make(chan event.DataEvent)
	event.EB.Subscribe("appRequestProgress", appCh)

	heartbeatTicker := time.NewTicker(30 * time.Second)

	ctx, cancel := context.WithCancel(r.Context())
	disconnect := ctx.Done()
	defer cancel()

	for {
		select {
		case appEvent := <-appCh:
			appRequestProgress, ok := appEvent.Data.(AppRequestProgress)
			if !ok {
				log.Println("Interface does not hold type App")
				return
			}
			err := sendAppToClient(appRequestProgress, clientManager)
			if err != nil {
				services.AppError(err.Error(), 500, w)
				return
			}
		case <-heartbeatTicker.C:
			err := sendMessage("heartbeat", userId, clientManager)
			if err != nil {
				log.Println("Error sending heartbeat: ", err)
				return
			}
		case <-disconnect:
			clientManager.RemoveClient(userId)
			event.EB.Unsubscribe("appRequestProgress", appCh)
			heartbeatTicker.Stop()
			cancel()
			log.Println("Client disconnected:", userId)
			return
		}
	}
}

func GetLiveRequestsRoute(router *mux.Router) {
	router.HandleFunc("/get-live", getLiveRequests).Methods("GET")
}
