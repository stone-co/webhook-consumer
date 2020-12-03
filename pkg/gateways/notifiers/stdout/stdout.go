package stdout

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/webhook-consumer/pkg/common/configuration"
	"github.com/stone-co/webhook-consumer/pkg/domain"
)

var _ domain.Notifier = &StdoutNotifier{}

type StdoutNotifier struct {
	log *logrus.Logger
}

func New() *StdoutNotifier {
	return &StdoutNotifier{}
}

func (n *StdoutNotifier) Configure(config *configuration.Config, log *logrus.Logger) error {
	n.log = log
	return nil
}

func (n StdoutNotifier) Send(ctx context.Context, eventTypeHeader, eventIDHeader, body string) error {
	log := n.log.WithField("notifier", "stdout")
	log.Printf("event headers: type[%s] id[%s]\n", eventTypeHeader, eventIDHeader)
	log.Printf("body: %s\n", body)

	return nil
}
