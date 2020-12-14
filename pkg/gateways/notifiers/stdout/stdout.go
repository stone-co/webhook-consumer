package stdout

import (
	"github.com/sirupsen/logrus"
	"github.com/stone-co/webhook-consumer/pkg/domain"
)

var _ domain.Notifier = &StdoutNotifier{}

type StdoutNotifier struct {
	log *logrus.Logger
}

func New() *StdoutNotifier {
	return &StdoutNotifier{}
}
