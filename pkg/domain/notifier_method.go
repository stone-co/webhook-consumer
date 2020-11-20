package domain

import (
	"context"
)

type NotifierMethod interface {
	Send(ctx context.Context, input CreateNotificationInput) error
}
