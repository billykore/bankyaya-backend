package messaging

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.bankyaya.org/app/backend/pkg/logger"
)

type consumeFunc func(context.Context) error

type Listener struct {
	log *logger.Logger
	tc  *TransferConsumer
}

func NewListener(log *logger.Logger, tc *TransferConsumer) *Listener {
	return &Listener{
		log: log,
		tc:  tc,
	}
}

func (l *Listener) Listen() {
	l.listen(map[string]consumeFunc{
		"transferConsumer": l.tc.Consume,
	})
}

func (l *Listener) listen(fn map[string]consumeFunc) {
	l.log.Info("Listening for messages")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for name, c := range fn {
		go func() {
			err := c(ctx)
			if err != nil {
				l.log.Errorf("error consuming %s: %v", name, err)
				cancel()
			}
		}()
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop // Wait for a termination signal
	l.log.Info("Shutting down listener...")

	time.Sleep(2 * time.Second)
	cancel()

	l.log.Info("Listener stopped")
	os.Exit(0)
}
