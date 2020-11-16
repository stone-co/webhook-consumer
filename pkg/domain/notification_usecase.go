package domain

import (
	"context"
)

type CreateNotificationInput struct {
	Body string
}

type NotificationUsecase interface {
	CreateNotification(ctx context.Context, input CreateNotificationInput) error
}
