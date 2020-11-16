package stdout

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/stone-co/webhook-consumer/pkg/domain"
)

var _ domain.NotifierMethod = &Stdout{}

type Stdout struct {
	log *logrus.Logger
}

func New(log *logrus.Logger) *Stdout {
	return &Stdout{
		log: log,
	}
}

func (std Stdout) Send(ctx context.Context, body string) error {
	std.log.Println(body)

	return nil
}
