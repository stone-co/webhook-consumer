package usecase

import (
	"github.com/sirupsen/logrus"

	"github.com/stone-co/webhook-consumer/pkg/domain"
)

var _ domain.NotificationUsecase = &NotificationUsecase{}

type NotificationUsecase struct {
	log     *logrus.Logger
	methods []domain.NotifierMethod
}

func NewNotificationUsecase(log *logrus.Logger, methods []domain.NotifierMethod) *NotificationUsecase {
	return &NotificationUsecase{
		log:     log,
		methods: methods,
	}
}
