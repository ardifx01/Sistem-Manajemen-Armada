package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

func StartConsumer(amqpURL string) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatal("gagal konek RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// 1. Buat queue
	q, err := ch.QueueDeclare(
		"geofence_alerts",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Bind queue ke exchange "fleet.events"
	err = ch.QueueBind(
		q.Name,
		"",
		"fleet.events",
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Listen pesan
	msgs, err := ch.Consume(
		q.Name,
		"",
		true, //auto ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(" [*] Menunggu pesan dari geofence_alerts ...")

	for msg := range msgs {
		log.Printf("📩 Pesan diterima: %s", msg.Body)
	}
}
