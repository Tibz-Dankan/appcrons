package request

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/middlewares"
	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/pubsub"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func sendMessage(w http.ResponseWriter, message, userId string) {
	data, _ := json.Marshal(map[string]string{
		"message": message,
		"userId":  userId,
	})

	w.Write([]byte("data: " + string(data) + "\n\n"))
}

func getLiveRequests(w http.ResponseWriter, r *http.Request) {
	log.Println("getting live request...")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	userId, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		services.AppError("UserID not found in context", 500, w)
		return
	}
	clientManager := services.NewClientManager()

	// Add the client to the manager
	clientManager.AddClient(userId, w)

	// Writing warmup message
	sendMessage(w, "warmup", userId)

	// Setting up interval for heartbeat message
	// ticker := time.NewTicker(30 * time.Second)
	// defer ticker.Stop()
	// for range ticker.C {
	// 	sendMessage(w, "heartbeat", userId)
	// }

	psub := pubsub.PubSub{}
	app := models.App{}

	userApps, err := app.FindByUser(userId)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	log.Println("About to tackle subscriptions")

	for _, app := range userApps {
		sub, err := psub.Subscribe(app.ID)
		if err != nil {
			log.Println("Error while getting pubsub:", err)
			services.AppError(err.Error(), 500, w)
		}

		client, ok := clientManager.GetClient(app.UserID)
		if !ok {
			log.Println("Client not found")
			return
		}

		subPayload, err := json.Marshal(sub)
		if err != nil {
			log.Println("Error converting sub payload to json:", err)
			services.AppError(err.Error(), 500, w)
		}

		_, err = client.Write([]byte("data: " + string(subPayload) + "\n\n"))
		if err != nil {
			log.Println("Error sending event:", err)
			services.AppError(err.Error(), 500, w)
		}
		client.(http.Flusher).Flush() // Flush the response writer to send data immediately
	}
}

func GetLiveRequestsRoute(router *mux.Router) {
	router.HandleFunc("/get-live", getLiveRequests).Methods("GET")
}
