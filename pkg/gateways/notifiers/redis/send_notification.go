package redis

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	RedisNotificationList = "STONE-NOTIFICATIONS"
)

type Notification struct {
	EventType string
	EventID   string
	Body      json.RawMessage
}

func (n RedisNotifier) Send(ctx context.Context, eventTypeHeader, eventIDHeader, body string) error {
	log := n.log.WithField("notifier", "redis")

	conn := n.pool.Get()
	defer conn.Close()

	notification := Notification{
		EventType: eventTypeHeader,
		EventID:   eventIDHeader,
		Body:      []byte(body),
	}

	encoded, err := json.Marshal(notification)
	if err != nil {
		log.WithError(err).Error("failed to marshal notification")
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	_, err = conn.Do("RPUSH", RedisNotificationList, encoded)
	if err != nil {
		log.WithError(err).Info("unable to save the notifier")
		return fmt.Errorf("unable to save the notifier: %w", err)
	}

	return nil
}
