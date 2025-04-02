package messaging

import (
	"context"

	"go.bankyaya.org/app/backend/internal/core/entity"
	"go.bankyaya.org/app/backend/internal/pkg/config"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/messaging/rabbitmq"
)

type SchedulerPublisher struct {
	cfg  *config.Config
	log  *logger.Logger
	conn *rabbitmq.Connection
}

func NewSchedulerPublisher(cfg *config.Config, log *logger.Logger, conn *rabbitmq.Connection) *SchedulerPublisher {
	return &SchedulerPublisher{
		cfg:  cfg,
		log:  log,
		conn: conn,
	}
}

func (sp *SchedulerPublisher) Publish(ctx context.Context, event *entity.TransferRequest) error {
	payload := rabbitmq.MessagePayload[*entity.TransferRequest]{
		Origin: "scheduler-service",
		Data:   event,
	}
	body, err := payload.MarshalBinary()
	if err != nil {
		sp.log.Errorf("MarshalBinary: %v", err)
		return err
	}

	err = sp.conn.Publish(ctx, sp.cfg.Rabbit.QueueName, body)
	if err != nil {
		sp.log.Errorf("Failed to publish message: %v", err)
		return err
	}

	sp.log.Infof("Published scheduler auto-debit. Schedule ID: %d", payload.Data.ScheduleId)
	return nil
}
