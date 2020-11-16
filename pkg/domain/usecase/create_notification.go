package usecase

import (
	"context"

	"github.com/stone-co/webhook-consumer/pkg/domain"
)

func (uc NotificationUsecase) CreateNotification(ctx context.Context, input domain.CreateNotificationInput) error {
	uc.method.Send(ctx, input.Body)
	return nil
}
