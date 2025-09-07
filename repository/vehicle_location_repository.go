package repository

import (
	"context"
	"log"
	"time"

	"transjakarta/model"

	"github.com/jackc/pgx/v5"
)

type VehicleLocationRepository interface {
	Insert(ctx context.Context, loc model.VehicleLocation) error
	GetLastLocation(ctx context.Context, vehicleID string) (*model.VehicleLocation, error)
	GetHistory(ctx context.Context, vehicleID string, start, end time.Time) ([]model.VehicleLocation, error)
}

type vehicleLocationRepository struct {
	db *pgx.Conn
}

func NewVehicleLocationRepository(db *pgx.Conn) VehicleLocationRepository {
	return &vehicleLocationRepository{db: db}
}

func (r *vehicleLocationRepository) Insert(ctx context.Context, loc model.VehicleLocation) error {
	log.Printf("Attempting to insert location for vehicle %s: lat=%f, lon=%f, ts=%v",
		loc.VehicleID, loc.Latitude, loc.Longitude, loc.Timestamp)
	_, err := r.db.Exec(ctx, `
        INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp) 
        VALUES ($1, $2, $3, $4)`,
		loc.VehicleID, loc.Latitude, loc.Longitude, loc.Timestamp,
	)
	if err != nil {
		log.Printf("Failed to insert location: %v", err)
	} else {
		log.Println(" Location inserted successfully")
	}
	return err
}

func (r *vehicleLocationRepository) GetLastLocation(ctx context.Context, vehicleID string) (*model.VehicleLocation, error) {
	row := r.db.QueryRow(ctx, `
		SELECT vehicle_id, latitude, longitude, timestamp 
		FROM vehicle_locations 
		WHERE vehicle_id = $1 
		ORDER BY timestamp DESC 
		LIMIT 1`, vehicleID)

	var loc model.VehicleLocation
	err := row.Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &loc.Timestamp)
	if err != nil {
		return nil, err
	}
	return &loc, nil
}

func (r *vehicleLocationRepository) GetHistory(ctx context.Context, vehicleID string, start, end time.Time) ([]model.VehicleLocation, error) {
	rows, err := r.db.Query(ctx, `
		SELECT vehicle_id, latitude, longitude, timestamp 
		FROM vehicle_locations 
		WHERE vehicle_id = $1 AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp ASC`, vehicleID, start, end)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []model.VehicleLocation
	for rows.Next() {
		var loc model.VehicleLocation
		if err := rows.Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &loc.Timestamp); err != nil {
			return nil, err
		}
		history = append(history, loc)
	}
	return history, nil
}
