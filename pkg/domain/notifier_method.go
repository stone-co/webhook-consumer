package domain

import (
	"context"

	"github.com/stone-co/webhook-consumer/pkg/common/configuration"
)

type Notifier interface {
	Configure(config *configuration.Config) error
	Send(ctx context.Context, input NotificationInput) error
}
