package domain

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Notifier interface {
	Configure(log *logrus.Logger) error
	Send(ctx context.Context, eventTypeHeader, eventIDHeader, body string) error
}
