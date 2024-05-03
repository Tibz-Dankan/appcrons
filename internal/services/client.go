package services

import (
	"net/http"
	"sync"
)

type ClientManager struct {
	clients map[string]http.ResponseWriter
	sync.RWMutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[string]http.ResponseWriter),
	}
}

func (cm *ClientManager) AddClient(userId string, w http.ResponseWriter) {
	cm.Lock()
	defer cm.Unlock()
	cm.clients[userId] = w
}

func (cm *ClientManager) RemoveClient(userId string) {
	cm.Lock()
	defer cm.Unlock()
	delete(cm.clients, userId)
}

func (cm *ClientManager) GetClient(userId string) (http.ResponseWriter, bool) {
	cm.RLock()
	defer cm.RUnlock()
	client, ok := cm.clients[userId]
	return client, ok
}
