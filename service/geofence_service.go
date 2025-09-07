package service

import (
	"math"

	"transjakarta/model"
	"transjakarta/rabbitmq"
)

type GeofenceService struct {
	publisher *rabbitmq.RabbitMQPublisher
	lat       float64
	lon       float64
	radiusM   float64
}

func NewGeofenceService(pub *rabbitmq.RabbitMQPublisher, lat, lon, radiusM float64) *GeofenceService {
	return &GeofenceService{publisher: pub, lat: lat, lon: lon, radiusM: radiusM}
}

func (g *GeofenceService) CheckAndPublish(loc model.VehicleLocation) error {
	distance := haversine(g.lat, g.lon, loc.Latitude, loc.Longitude)
	if distance <= g.radiusM {
		event := map[string]interface{}{
			"vehicle_id": loc.VehicleID,
			"event":      "geofence_entry",
			"location": map[string]float64{
				"latitude":  loc.Latitude,
				"longitude": loc.Longitude,
			},
			"timestamp": loc.Timestamp.Unix(),
		}
		return g.publisher.Publish(event)
	}
	return nil
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
