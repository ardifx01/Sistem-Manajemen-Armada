package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	broker := "tcp://localhost:1883"
	clientID := "mock-publisher"
	vehicleID := "B1234XYZ"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("gagal konek MQTT broker: %v", token.Error())
	}
	defer client.Disconnect(250)

	log.Printf(" Publisher terhubung ke %s", broker)

	baseLat := -6.2088
	baseLon := 106.8456

	for {
		lat := baseLat + (rand.Float64()-0.5)/1000
		lon := baseLon + (rand.Float64()-0.5)/1000
		ts := time.Now().Unix()

		payload := map[string]interface{}{
			"vehicle_id": vehicleID,
			"latitude":   lat,
			"longitude":  lon,
			"timestamp":  ts,
		}

		data, _ := json.Marshal(payload)
		topic := fmt.Sprintf("/fleet/vehicle/%s/location", vehicleID)

		client.Publish(topic, 1, false, data).Wait()
		log.Printf("Published ke %s: %s", topic, string(data))

		time.Sleep(2 * time.Second)
	}
}
