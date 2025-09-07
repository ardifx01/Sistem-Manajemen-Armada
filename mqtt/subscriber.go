package mqtt

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"transjakarta/model"
	"transjakarta/service"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTSubscriber struct {
	client   mqtt.Client
	service  service.VehicleLocationService
	geofence *service.GeofenceService
}

func NewMQTTSubscriber(broker, clientID string, svc service.VehicleLocationService, gf *service.GeofenceService) *MQTTSubscriber {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)

	client := mqtt.NewClient(opts)
	// if token := client.Connect(); token.Wait() && token.Error() != nil {
	// 	log.Fatalf("gagal konek MQTT broker: %v", token.Error())
	// }

	for i := 0; i < 10; i++ {
		token := client.Connect()
		if token.Wait() && token.Error() == nil {
			log.Println(" Connected to MQTT broker")
			break
		}
		log.Printf("gagal konek MQTT broker: %v", token.Error())
		time.Sleep(3 * time.Second)
	}

	return &MQTTSubscriber{client: client, service: svc, geofence: gf}
}

func (s *MQTTSubscriber) Subscribe() {
	topic := "/fleet/vehicle/+/location"
	if token := s.client.Subscribe(topic, 1, s.handleMessage); token.Wait() && token.Error() != nil {
		log.Fatalf("gagal subscribe: %v", token.Error())
	}
	log.Printf(" Subscriber aktif di topic %s", topic)
}

func (s *MQTTSubscriber) handleMessage(client mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	if len(parts) < 4 {
		log.Println("topic format salah")
		return
	}

	var payload struct {
		VehicleID string  `json:"vehicle_id"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timestamp int64   `json:"timestamp"`
	}
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Printf("gagal decode JSON: %v", err)
		return
	}

	loc := model.VehicleLocation{
		VehicleID: payload.VehicleID, // Gunakan vehicle_id dari payload
		Latitude:  payload.Latitude,
		Longitude: payload.Longitude,
		Timestamp: time.Unix(payload.Timestamp, 0),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := s.service.SaveLocation(ctx, loc); err != nil {
		log.Printf("gagal simpan DB: %v", err)
		return
	}
	// log.Printf(" data lokasi %s tersimpan", vehicleID)
	if s.geofence != nil {
		if err := s.geofence.CheckAndPublish(loc); err != nil {
			log.Printf("gagal publish geofence event: %v", err)
		}
	}
}
