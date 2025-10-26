package metrics

import "github.com/sirupsen/logrus"

type NoOp struct{}

func (d NoOp) Inc(_ string, _ int) {
	// todo - add some implementation!
	logrus.Info("NoOp metrics client")
}
