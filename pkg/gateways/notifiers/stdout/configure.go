package stdout

import (
	"github.com/sirupsen/logrus"
)

func (n *StdoutNotifier) Configure(log *logrus.Logger) error {
	n.log = log

	return nil
}
