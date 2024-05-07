package request

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/event"
	"github.com/Tibz-Dankan/keep-active/internal/middlewares"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func sendMessage(w http.ResponseWriter, message, userId string) {
	data, _ := json.Marshal(map[string]string{
		"message": message,
		"userId":  userId,
	})

	w.Write([]byte("data: " + string(data) + "\n\n"))
	w.(http.Flusher).Flush()
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

	_, err = client.Write([]byte("data: " + string(appJson) + "\n\n"))
	if err != nil {
		log.Println("Error sending event:", err)
		return err
	}
	client.(http.Flusher).Flush()

	return nil
}

func getLiveRequests(w http.ResponseWriter, r *http.Request) {
	log.Println("getting live request...")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.(http.Flusher).Flush()

	userId, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		services.AppError("UserID not found in context", 500, w)
		return
	}
	clientManager := services.NewClientManager()
	clientManager.AddClient(userId, w)

	// Writing warmup message
	sendMessage(w, "warmup", userId)
	// TODO: send message to client every 30 seconds

	appCh := make(chan event.DataEvent)

	event.EB.Subscribe("app", appCh)

	type App = models.App

	// Listening for events
	for {
		appEvent := <-appCh

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
	}
}

func GetLiveRequestsRoute(router *mux.Router) {
	router.HandleFunc("/get-live", getLiveRequests).Methods("GET")
}
