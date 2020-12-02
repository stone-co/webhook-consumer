package stdout

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/webhook-consumer/pkg/common/configuration"
	"github.com/stone-co/webhook-consumer/pkg/domain"
)

var _ domain.NotifierMethod = &StdoutNotifier{}

type StdoutNotifier struct {
	log *logrus.Logger
}

func New(log *logrus.Logger) *StdoutNotifier {
	return &StdoutNotifier{
		log: log,
	}
}

func (n StdoutNotifier) Configure(config *configuration.Config) error {
	return nil
}

func (n StdoutNotifier) Send(ctx context.Context, input domain.NotificationInput) error {
	log := n.log.WithField("notifier", "stdout")
	log.Printf("Body: %s\n", input.Body)
	log.Printf("Header: %+v\n", input.Header)

	return nil
}
