package notifications

import (
	"github.com/sirupsen/logrus"

	"github.com/stone-co/webhook-consumer/pkg/common/validator"
	"github.com/stone-co/webhook-consumer/pkg/domain"
)

type Handler struct {
	log *logrus.Logger
	*validator.JSONValidator
	privateKey      interface{}
	verificationKey interface{}
	usecase         domain.NotificationUsecase
}

func NewHandler(log *logrus.Logger, validator *validator.JSONValidator, privateKey interface{}, verificationKey interface{}, usecase domain.NotificationUsecase) *Handler {
	return &Handler{
		log:             log,
		JSONValidator:   validator,
		privateKey:      privateKey,
		verificationKey: verificationKey,
		usecase:         usecase,
	}
}
