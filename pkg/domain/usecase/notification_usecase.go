package usecase

import (
	"github.com/sirupsen/logrus"

	"github.com/stone-co/webhook-consumer/pkg/domain"
)

var _ domain.NotificationUsecase = &NotificationUsecase{}

type NotificationUsecase struct {
	log       *logrus.Logger
	notifiers []domain.Notifier
}

func NewNotificationUsecase(log *logrus.Logger, notifiers []domain.Notifier) *NotificationUsecase {
	return &NotificationUsecase{
		log:       log,
		notifiers: notifiers,
	}
}
