package event

import (
	"log"
	"sync"
)

type DataEvent struct {
	Data  interface{}
	Topic string
}

type DataChannel chan DataEvent

type DataChannelSlice []DataChannel

type EventBus struct {
	subscribers map[string]DataChannelSlice
	rm          sync.RWMutex
}

func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.rm.RLock()
	log.Println("received publish data:::", data)

	if chans, found := eb.subscribers[topic]; found {

		channels := append(DataChannelSlice{}, chans...)
		go func(data DataEvent, dataChannelSlices DataChannelSlice) {
			for _, ch := range dataChannelSlices {
				ch <- data
			}
		}(DataEvent{Data: data, Topic: topic}, channels)
	}
	eb.rm.RUnlock()
}

func (eb *EventBus) Subscribe(topic string, ch DataChannel) {
	eb.rm.Lock()
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
		log.Println("subscribed to channel :::", ch)

	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
		log.Println("subscribed to channel :::", ch)

	}
	eb.rm.Unlock()
}

var EB = &EventBus{
	subscribers: map[string]DataChannelSlice{},
}
