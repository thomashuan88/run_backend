package queue

import (
	"fmt"
	// "os"

	"github.com/streadway/amqp"
)

var RabbitMqClient = Conn{}

// Conn -
type Conn struct {
	Channel *amqp.Channel
}

// 初始化rabbitMq连接
func RabbitMq(url string) {
	conn, err := GetConn(url)
	if err != nil {
		panic("rabbitmq连接异常: " + err.Error())
	}

	RabbitMqClient = conn
}

// GetConn -
func GetConn(rabbitURL string) (Conn, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return Conn{}, err
	}

	ch, err := conn.Channel()
	return Conn{
		Channel: ch,
	}, err
}

// Publish -
func (conn Conn) Publish(routingKey string, data []byte) error {
	return conn.Channel.Publish(
		// exchange - yours may be different
		"amq.topic",
		routingKey,
		// mandatory - we don't care if there I no queue
		false,
		// immediate - we don't care if there is no consumer on the queue
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent,
		})
}

// StartConsumer -
func (conn Conn) StartConsumer(
	queueName,
	routingKey string,
	handler func(d amqp.Delivery) bool,
	concurrency int) error {

	// create the queue if it doesn't already exist
	_, err := conn.Channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// bind the queue to the routing key
	err = conn.Channel.QueueBind(queueName, routingKey, "amq.topic", false, nil)
	if err != nil {
		return err
	}

	// prefetch 4x as many messages as we can handle at once
	prefetchCount := concurrency * 4
	err = conn.Channel.Qos(prefetchCount, 0, false)
	if err != nil {
		return err
	}

	msgs, err := conn.Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	// create a goroutine for the number of concurrent threads requested
	for i := 0; i < concurrency; i++ {
		fmt.Printf("Processing messages on thread %v...\n", i)
		// go func() {

		for msg := range msgs {
			fmt.Println(msg)
			// if tha handler returns true then ACK, else NACK
			// the message back into the rabbit queue for
			// another round of processing
			if handler(msg) {
				msg.Ack(false)
			} else {
				msg.Nack(false, true)
			}
		}
		fmt.Println("Rabbit consumer closed - critical Error")
		// os.Exit(1)
		// }()
	}
	return nil
}
