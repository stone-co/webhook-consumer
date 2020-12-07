package domain

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/webhook-consumer/pkg/common/configuration"
)

type Notifier interface {
	Configure(config *configuration.Config, log *logrus.Logger) error
	Send(ctx context.Context, eventTypeHeader, eventIDHeader, body string) error
}
