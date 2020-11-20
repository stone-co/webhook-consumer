package usecase

import (
	"context"

	"github.com/stone-co/webhook-consumer/pkg/domain"
)

func (uc NotificationUsecase) SendNotification(ctx context.Context, input domain.NotificationInput) error {
	err := uc.method.Send(ctx, input)
	return err
}
