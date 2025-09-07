package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"transjakarta/config"
	"transjakarta/handler"
	"transjakarta/mqtt"
	"transjakarta/rabbitmq"
	"transjakarta/repository"
	"transjakarta/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbConn, err := config.InitPostgres()
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	// Init repository & service
	repo := repository.NewVehicleLocationRepository(dbConn)
	locService := service.NewVehicleLocationService(repo)

	rmqURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	rmq := rabbitmq.NewRabbitMQPublisher(rmqURL, "fleet.events")
	geofence := service.NewGeofenceService(rmq, -6.2088, 106.8456, 50) // contoh titik Monas
	go rabbitmq.StartConsumer(rmqURL)

	// Init MQTT Subscriber
	broker := fmt.Sprintf("tcp://%s:%s",
		os.Getenv("MQTT_BROKER"),
		os.Getenv("MQTT_PORT"),
	)

	subscriber := mqtt.NewMQTTSubscriber(broker, "fleet-subscriber", locService, geofence)
	subscriber.Subscribe()

	// Init Fiber
	app := fiber.New()
	// API Routes
	vh := handler.NewVehicleHandler(locService)
	app.Get("/vehicles/:vehicle_id/location", vh.GetLastLocation)
	app.Get("/vehicles/:vehicle_id/history", vh.GetHistory)

	log.Println("🚀 Server running on :8090")
	if err := app.Listen(":8090"); err != nil {
		log.Fatal(err)
	}
}
