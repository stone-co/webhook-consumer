package usecase

import (
	"github.com/sirupsen/logrus"

	"github.com/stone-co/webhook-consumer/pkg/common/keys"
	"github.com/stone-co/webhook-consumer/pkg/domain"
)

var _ domain.NotificationUsecase = &NotificationUsecase{}

type NotificationUsecase struct {
	log       *logrus.Logger
	keys      *keys.Config
	notifiers []domain.Notifier
}

func NewNotificationUsecase(log *logrus.Logger, keys *keys.Config, notifiers []domain.Notifier) *NotificationUsecase {
	return &NotificationUsecase{
		log:       log,
		keys:      keys,
		notifiers: notifiers,
	}
}
