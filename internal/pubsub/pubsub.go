package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Tibz-Dankan/keep-active/internal/config"
)

var redisClient = config.RedisClient()
var ctx = context.Background()

type PubSub struct{}

func (ps *PubSub) Publish(channel string, payload interface{}) error {

	// convert payload to json
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	if err := redisClient.Publish(ctx, channel, jsonPayload).Err(); err != nil {
		return err
	}

	return nil
}

func (ps *PubSub) Subscribe(channel string) (interface{}, error) {

	pubsub := redisClient.Subscribe(ctx, channel)
	var payload interface{}

	// Wait for messages
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			return payload, err
		}

		// Unmarshal the JSON payload into the payload
		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			return payload, err
		}

		fmt.Println(msg.Channel, payload)

		if payload != nil {
			return payload, nil
		}
	}

}
