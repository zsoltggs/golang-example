package notifier

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zsoltggs/golang-example/services/users/pkg/userevents"
)

type Notifier interface {
	Notify(ctx context.Context, event userevents.UserChangedEvent) error
}

type logNotifier struct {
	// Possibility to store generic logger interface here
}

func NewLogNotifier() Notifier {
	return &logNotifier{}
}

func (l logNotifier) Notify(_ context.Context, event userevents.UserChangedEvent) error {
	logrus.Infof("user changed with id %q", event.UserID)
	return nil
}
