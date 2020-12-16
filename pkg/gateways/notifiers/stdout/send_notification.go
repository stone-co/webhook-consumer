package stdout

import (
	"context"
)

func (n StdoutNotifier) Send(ctx context.Context, eventTypeHeader, eventIDHeader, body string) error {
	log := n.log.WithField("notifier", "stdout")
	log.Printf("event headers: type[%s] id[%s]\n", eventTypeHeader, eventIDHeader)
	log.Printf("body: %s\n", body)

	return nil
}
