package usecase

import (
	"context"

	"github.com/stone-co/webhook-consumer/pkg/domain"
)

func (uc NotificationUsecase) CreateNotification(ctx context.Context, input domain.CreateNotificationInput) error {
	err := uc.method.Send(ctx, input.Body)
	return err
}
