package domain

import (
	"context"
)

type CreateNotificationInput struct {
	Header HeaderNotification
	Body   string
}

type HeaderNotification struct {
	EventID   string
	EventType string
}

type NotificationUsecase interface {
	CreateNotification(ctx context.Context, input CreateNotificationInput) error
}
