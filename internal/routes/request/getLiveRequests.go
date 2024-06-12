package request

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Tibz-Dankan/keep-active/internal/event"
	"github.com/Tibz-Dankan/keep-active/internal/middlewares"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func sendMessage(w http.ResponseWriter, message, userId string) {
	f, ok := w.(http.Flusher)
	if !ok {
		log.Println("Response writer does not implement http.Flusher")
		return
	}

	data, _ := json.Marshal(map[string]string{
		"message": message,
		"userId":  userId,
	})

	_, err := w.Write([]byte("data: " + string(data) + "\n\n"))
	if err != nil {
		log.Println("Error writing to response writer:", err)
		return
	}
	f.Flush()
}

func sendAppToClient(app models.App, clientManager *services.ClientManager) error {
	client, ok := clientManager.GetClient(app.UserID)
	if !ok {
		log.Println("Client not found")
		return nil
	}

	appJson, err := json.Marshal(&app)
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
	clientManager := services.NewClientManager()
	clientManager.AddClient(userId, w)

	// Writing warmup message
	sendMessage(w, "warmup", userId)

	appCh := make(chan event.DataEvent)
	// defer close(appCh)

	event.EB.Subscribe("app", appCh)

	type App = models.App

	heartbeatTicker := time.NewTicker(30 * time.Second)
	defer heartbeatTicker.Stop()

	// Listening for events and sending heartbeats
	for {
		select {
		case appEvent := <-appCh:
			app, ok := appEvent.Data.(App)
			if !ok {
				log.Println("Interface does not hold type App")
				return
			}
			err := sendAppToClient(app, clientManager)
			if err != nil {
				services.AppError(err.Error(), 500, w)
				return
			}
		default:
			select {
			case <-heartbeatTicker.C:
				sendMessage(w, "heartbeat", userId)
				// default:
			}
		}
	}
}

func GetLiveRequestsRoute(router *mux.Router) {
	router.HandleFunc("/get-live", getLiveRequests).Methods("GET")
}
