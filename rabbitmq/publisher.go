package rabbitmq

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQPublisher struct {
	channel  *amqp.Channel
	exchange string
}

// NewRabbitMQPublisher dengan retry agar tidak langsung fatal kalau RabbitMQ belum ready
func NewRabbitMQPublisher(amqpURL, exchange string) *RabbitMQPublisher {
	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error

	for i := 1; i <= 10; i++ { // coba 10x dengan delay
		conn, err = amqp.Dial(amqpURL)
		if err == nil {
			ch, err = conn.Channel()
			if err == nil {
				err = ch.ExchangeDeclare(exchange, "fanout", true, false, false, false, nil)
				if err == nil {
					log.Println(" Connected to RabbitMQ & exchange declared:", exchange)
					return &RabbitMQPublisher{channel: ch, exchange: exchange}
				}
			}
		}

		log.Printf("gagal konek RabbitMQ (try %d/10): %v", i, err)
		time.Sleep(5 * time.Second)
	}

	log.Fatal("tidak bisa konek RabbitMQ setelah 10 percobaan")
	return nil
}

func (p *RabbitMQPublisher) Publish(event interface{}) error {
	body, _ := json.Marshal(event)
	return p.channel.Publish(
		p.exchange,
		"", // fanout: routing key kosong
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
