package queue_manager

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
	"github.com/PAY-HERO-CONSULTING/gh-tools/dtos"
	"github.com/PAY-HERO-CONSULTING/gh-tools/internal/brokers"
	"github.com/PAY-HERO-CONSULTING/gh-tools/utils"
)

type (
	QueueManager interface {
		// rabbitmq
		Close() error
		Enqueue(ctx context.Context, queueName string, data any) error
		EnqueueWithDelay(ctx context.Context, queueName string, data any, delay time.Duration) error
		SubscribeToQueue(ctx context.Context, wg *sync.WaitGroup, queueName string, consumerFunc func(*dtos.Delivery)) error
	}

	queueManager struct {
		messageBroker brokers.Broker
	}
)

func NewQueueManager(
	kafkaHost,
	kafkaPort,
	brokerUrl string,
) (QueueManager, error) {
	messageBroker := brokers.NewMessageBroker(
		brokerUrl,
	)

	err := messageBroker.Connect()
	if err != nil {
		return nil, err
	}

	return &queueManager{
		messageBroker: messageBroker,
	}, nil
}

func (m *queueManager) Close() error {
	if m.messageBroker != nil {
		err := m.messageBroker.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *queueManager) Enqueue(ctx context.Context, queueName string, data any) error {
	body, err := m.marshal(data)
	if err != nil {
		return apperrs.New(
			err,
			apperrs.Internal,
		).LogErrorMessage(
			"json marshal request [%+v]",
			data,
		)
	}

	return m.messageBroker.PublishToQueue(ctx, queueName, body)
}

func (m *queueManager) EnqueueWithDelay(ctx context.Context, queueName string, data any, delay time.Duration) error {
	body, err := m.marshal(data)
	if err != nil {
		return apperrs.New(
			err,
			apperrs.Internal,
		).LogErrorMessage(
			"json marshal request delay [%+v]",
			data,
		)
	}

	return m.messageBroker.PublishToQueueWithBackOff(ctx, queueName, body, delay)
}

func (m *queueManager) marshal(data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (m *queueManager) SubscribeToQueue(
	ctx context.Context,
	wg *sync.WaitGroup,
	queueName string,
	consumerFunc func(*dtos.Delivery),
) error {
	return m.messageBroker.SubscribeToQueue(ctx, queueName, utils.ExtractConsumer(queueName), wg, consumerFunc)
}
