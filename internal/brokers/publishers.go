package brokers

import (
	"context"
	"log"
	"time"

	"github.com/PAY-HERO-CONSULTING/gh-tools/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (m *rabbitMQBroker) PublishToQueue(ctx context.Context, queue string, body []byte) error {
	return m.publishToQueue(ctx, queue, "", false, body, make(amqp.Table))
}

// Connecting to queues
// what is the queue
// what is the exhange
func (m *rabbitMQBroker) publishToQueue(
	ctx context.Context,
	queue,
	exchange string,
	delayed bool,
	body []byte,
	headers amqp.Table,
) error {
	if m.conn == nil {
		return errConnNotInitialized
	}

	ch, err := m.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	_, err = m.declareQueue(ch, queue, exchange, delayed)
	if err != nil {
		return err
	}

	var mandatory, immediate bool

	msg := amqp.Publishing{
		Body:         body,
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Headers:      headers,
		Timestamp:    time.Now(),
	}

	return ch.PublishWithContext(ctx, exchange, queue, mandatory, immediate, msg)
}

func (m *rabbitMQBroker) PublishToQueueWithBackOff(
	ctx context.Context,
	queue string,
	body []byte,
	delay time.Duration,
) error {
	delayMillis := delay.Nanoseconds() / 1e6

	headers := make(amqp.Table)
	headers["x-delay"] = delayMillis

	return m.publishToQueue(ctx, queue, delayedExchangeName, true, body, headers)
}

func (m *rabbitMQBroker) declareQueue(ch *amqp.Channel, queue, exchange string, delayed bool) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		queue, // name of the queue
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		log.Printf("Failed to register queue: %v", err)
		return q, err
	}

	if delayed {
		err = ch.QueueBind(q.Name, q.Name, exchange, false, nil)
		if err != nil {
			log.Printf("Failed to bind queue: %v", err)
			return q, err
		}
	}

	return q, nil
}

func (m *rabbitMQBroker) PublishToTopic(
	ctx context.Context,
	exchange string,
	notificationType string,
	body []byte,
) error {
	ch, err := m.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	logger.Infof("Pushing event: [%+v]", string(body))

	err = ch.PublishWithContext(
		ctx,
		exchange,         // exchange
		notificationType, // routing key
		false,            // mandatory?
		false,            // immediate?
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		logger.Infof("failed to publish topic err: [%+v]", err)
		return err
	}

	logger.Info("sucessfully published a topic")

	return nil
}
