package brokers

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/PAY-HERO-CONSULTING/gh-tools/dtos"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker interface {
	Close() error
	Connect() error
	PublishToQueue(ctx context.Context, queue string, body []byte) error
	PublishToQueueWithBackOff(ctx context.Context, queue string, body []byte, delay time.Duration) error
	PublishToTopic(ctx context.Context, exchange string, notificationType string, body []byte) error
	Shutdown()
	SubscribeToBackedOffQueue(ctx context.Context, queue, consumer string, wg *sync.WaitGroup, consumerFunc func(*dtos.Delivery)) error
	SubscribeToQueue(ctx context.Context, queue, consumer string, wg *sync.WaitGroup, consumerFunc func(*dtos.Delivery)) error
}

type rabbitMQBroker struct {
	brokerUrl string
	closeChan chan *amqp.Error
	conn      *amqp.Connection
	quitChan  chan bool
}

func NewMessageBroker(
	url string,
) Broker {
	return NewMessageBrokerWithParams(url)
}

func NewMessageBrokerWithParams(
	brokerUrl string,
) Broker {
	return &rabbitMQBroker{
		brokerUrl: brokerUrl,
		closeChan: make(chan *amqp.Error),
		quitChan:  make(chan bool),
	}
}

func (m *rabbitMQBroker) Close() error {
	if m.conn != nil {
		return m.conn.Close()
	}

	return nil
}

func (m *rabbitMQBroker) Connect() error {
	err := m.attemptToConnect()
	if err != nil {
		return err
	}

	err = m.declareDirectExchange()
	if err != nil {
		return err
	}

	err = m.declareDelayedExchange()
	if err != nil {
		return err
	}

	go m.handleDisconnect()

	return nil
}

func (m *rabbitMQBroker) attemptToConnect() error {
	var err error
	m.conn, err = amqp.Dial(m.brokerUrl)
	if err != nil {
		return err
	}

	m.conn.NotifyClose(m.closeChan)
	return nil
}

func (m *rabbitMQBroker) declareDirectExchange() error {
	return m.declareExchange(directExchangeName, amqp.ExchangeDirect, nil)
}

func (m *rabbitMQBroker) declareDelayedExchange() error {
	args := make(amqp.Table)
	args["x-delayed-type"] = amqp.ExchangeDirect

	return m.declareExchange(delayedExchangeName, "x-delayed-message", args)
}

// Shutdown closes rabbitmq's connection
func (m *rabbitMQBroker) Shutdown() {
	m.quitChan <- true
	log.Println("shutting down rabbitMQ's connection...")
}

func (m *rabbitMQBroker) declareExchange(name, kind string, args amqp.Table) error {
	ch, err := m.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	return ch.ExchangeDeclare(
		name,  // name of the exchange
		kind,  // type
		true,  // durable
		false, // auto-deleted?
		false, // internal
		false, // no-wait?
		args,
	)
}

// handleDisconnect will try to reconnect every 5 seconds after a disconnection
func (m *rabbitMQBroker) handleDisconnect() {
	for {
		select {
		case errChan := <-m.closeChan:
			if errChan != nil {
				log.Printf("rabbitMQ disconnection: %v", errChan)
			}
		case <-m.quitChan:
			err := m.Close()
			if err != nil {
				log.Printf("rabbitMQ close error: %v", err)
			}

			log.Println("...rabbitMQ has been shut down")
			m.quitChan <- true
			return
		}

		time.Sleep(5 * time.Second)
		log.Println("trying to reconnect to rabbitMQ...")

		err := m.attemptToConnect()
		if err != nil {
			log.Fatalf("failed to attempt reconnection error: %v", err)
		}
	}
}

// func (r *RabbitMQBroker) Close(ctx context.Context) (done chan struct{}) {

// 	done = make(chan struct{})

// 	doneWaiting := make(chan struct{})
// 	go func() {
// 		r.wg.Wait()
// 		close(doneWaiting)
// 	}()

// 	go func() {
// 		defer close(done)
// 		select { // either waits for the messages to process or timeout from context
// 		case <-doneWaiting:
// 		case <-ctx.Done():
// 		}
// 		closeConnections(r)
// 	}()
// 	return
// }
