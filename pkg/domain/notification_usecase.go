package domain

import (
	"context"
)

type NotificationInput struct {
	Header        HeaderNotification
	EncryptedBody string
}

type HeaderNotification struct {
	EventID   string
	EventType string
}

type NotificationUsecase interface {
	SendNotification(ctx context.Context, input NotificationInput) error
}
