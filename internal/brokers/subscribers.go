package brokers

import (
	"context"
	"log"
	"sync"

	"github.com/PAY-HERO-CONSULTING/gh-tools/dtos"
	"github.com/PAY-HERO-CONSULTING/gh-tools/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (m *rabbitMQBroker) SubscribeToBackedOffQueue(ctx context.Context, queue, consumer string, wg *sync.WaitGroup, consumerFunc func(*dtos.Delivery)) error {
	return m.subscribeToQueue(ctx, queue, consumer, delayedExchangeName, true, wg, consumerFunc)
}

func (m *rabbitMQBroker) SubscribeToQueue(ctx context.Context, queue, consumer string, wg *sync.WaitGroup, consumerFunc func(*dtos.Delivery)) error {
	return m.subscribeToQueue(ctx, queue, consumer, "", false, wg, consumerFunc)
}

func (m *rabbitMQBroker) subscribeToQueue(
	ctx context.Context,
	queue,
	consumer,
	exchange string,
	delayed bool,
	wg *sync.WaitGroup,
	consumerFunc func(*dtos.Delivery),
) error {
	if m.conn == nil {
		return errConnNotInitialized
	}

	ch, err := m.conn.Channel()
	if err != nil {
		return err
	}

	logger.Infof("Declaring Queue (%s)", queue)

	q, err := m.declareQueue(ch, queue, exchange, delayed)
	if err != nil {
		return err
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Printf("Failed on Qos: %v", err)
		return err
	}

	deliveryChan, err := ch.Consume(
		q.Name,
		consumer,
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
		return err
	}

	go m.consumeDeliveries(ctx, queue, wg, deliveryChan, consumerFunc)

	return nil
}

func (m *rabbitMQBroker) consumeDeliveries(
	ctx context.Context,
	queue string,
	wg *sync.WaitGroup,
	deliveryChan <-chan amqp.Delivery,
	consumerFunc func(d *dtos.Delivery),
) {
	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case delivery := <-deliveryChan:
			d := &dtos.Delivery{
				Body:            delivery.Body,
				ContentEncoding: delivery.ContentEncoding,
				ContentType:     delivery.ContentType,
				Exchange:        delivery.Exchange,
				Expiration:      delivery.Expiration,
				Headers:         delivery.Headers,
				MessageId:       delivery.MessageId,
				RoutingKey:      delivery.RoutingKey,
				Timestamp:       delivery.Timestamp,
				Type:            delivery.Type,
			}

			consumerFunc(d)

			delivery.Ack(false)
		case <-ctx.Done():
			logger.Infof("Shut down consumer for [%v]", queue)
			return
		}
	}
}

func (m *rabbitMQBroker) declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name
		false, // durable?
		false, // delete when unsued?
		true,  // exclusive?
		false, // no-wait?
		nil,   // arguements?
	)
}

func (m *rabbitMQBroker) consumeTopics(
	ctx context.Context,
	topicsExhange,
	topic string,
	deliveryChan <-chan amqp.Delivery,
	consumerFunc func(d *dtos.Delivery),
) {
	// creating a channel to consume RabbitMQ messages forever
	forever := make(chan bool)

	go func() {
		for delivery := range deliveryChan {
			payload := &dtos.Delivery{
				Body:            delivery.Body,
				ContentEncoding: delivery.ContentEncoding,
				ContentType:     delivery.ContentType,
				Exchange:        delivery.Exchange,
				Expiration:      delivery.Expiration,
				Headers:         delivery.Headers,
				MessageId:       delivery.MessageId,
				RoutingKey:      delivery.RoutingKey,
				Timestamp:       delivery.Timestamp,
				Type:            delivery.Type,
			}

			consumerFunc(payload)
			delivery.Ack(false)
		}
	}()

	logger.Infof(" [*] Waiting for messages on [Exchange, Queue] [%s, %s].\n", topicsExhange, topic)
	<-forever
}
