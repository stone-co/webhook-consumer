package usecase

import (
	"context"

	"github.com/stone-co/webhook-consumer/pkg/domain"
)

func (uc NotificationUsecase) SendNotification(ctx context.Context, input domain.NotificationInput) error {
	for _, method := range uc.methods {
		err := method.Send(ctx, input)
		if err != nil {
			return err
		}
	}

	return nil
}
