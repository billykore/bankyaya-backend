package messaging

import (
	"context"

	"go.bankyaya.org/app/backend/pkg/entity"
	"go.bankyaya.org/app/backend/pkg/service"
	"go.bankyaya.org/app/backend/pkg/util/config"
	"go.bankyaya.org/app/backend/pkg/util/logger"
	rabbitmq "go.bankyaya.org/app/backend/pkg/util/messaging/rabbitmq"
)

type TransferConsumer struct {
	cfg  *config.Config
	log  *logger.Logger
	conn *rabbitmq.Connection
	svc  *service.Transfer
}

func NewTransferConsumer(cfg *config.Config, log *logger.Logger, conn *rabbitmq.Connection, svc *service.Transfer) *TransferConsumer {
	return &TransferConsumer{
		cfg:  cfg,
		log:  log,
		conn: conn,
		svc:  svc,
	}
}

func (tc *TransferConsumer) Consume(ctx context.Context) error {
	msgs, err := tc.conn.Consume(ctx, tc.cfg.Rabbit.QueueName)
	if err != nil {
		tc.log.Errorf("Consume returns error: %v", err)
		return err
	}

	for msg := range msgs {
		go func() {
			tc.log.Infof("Received a message: %s", msg.Body)

			payload := new(rabbitmq.MessagePayload[entity.Event])
			err := payload.UnmarshalBinary(msg.Body)
			if err != nil {
				tc.log.Infof("Failed to process user event: %v", err)
				return
			}

			transaction, err := tc.svc.ProcessAutoTransferEvent(ctx, &payload.Data)
			if err != nil {
				tc.log.Errorf("Failed to process transfer event: %v", err)
			}

			tc.log.Infof("Process success. Data: %v", transaction)
		}()
	}

	return nil
}
