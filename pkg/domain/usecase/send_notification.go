package usecase

import (
	"context"

	"github.com/stone-co/webhook-consumer/pkg/domain"
)

func (uc NotificationUsecase) SendNotification(ctx context.Context, input domain.NotificationInput) error {
	for _, notifier := range uc.notifiers {
		err := notifier.Send(ctx, input)
		if err != nil {
			return err
		}
	}

	return nil
}
