package request

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/event"
	"github.com/Tibz-Dankan/keep-active/internal/middlewares"
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
	clientManager.AddClient(userId, w)

	// Writing warmup message
	sendMessage(w, "warmup", userId)

	appCh := make(chan event.DataEvent)

	event.EB.Subscribe("app", appCh)

	// type App = models.App

	// Start listening for events
	go func() {
		for {
			// event := <-listener
			app := <-appCh
			log.Println("Received subscription  app data from event:::", app)
			// var userApp App  = App{app.Data}
			// send request to the client
		}
	}()
	// 	client, ok := clientManager.GetClient(app.UserID)
	// 	if !ok {
	// 		log.Println("Client not found")
	// 		return
	// 	}

	// 	_, err = client.Write([]byte("data: " + string(subPayload) + "\n\n"))
	// 	if err != nil {
	// 		log.Println("Error sending event:", err)
	// 		services.AppError(err.Error(), 500, w)
	// 	}
	// 	client.(http.Flusher).Flush()

}

func GetLiveRequestsRoute(router *mux.Router) {
	router.HandleFunc("/get-live", getLiveRequests).Methods("GET")
}
