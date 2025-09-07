package service

import (
	"context"
	"time"

	"transjakarta/model"
	"transjakarta/repository"
)

type VehicleLocationService interface {
	SaveLocation(ctx context.Context, loc model.VehicleLocation) error
	GetLastLocation(ctx context.Context, vehicleID string) (*model.VehicleLocation, error)
	GetHistory(ctx context.Context, vehicleID string, start, end time.Time) ([]model.VehicleLocation, error)
}

type vehicleLocationService struct {
	repo repository.VehicleLocationRepository
}

func NewVehicleLocationService(repo repository.VehicleLocationRepository) VehicleLocationService {
	return &vehicleLocationService{repo: repo}
}

func (s *vehicleLocationService) SaveLocation(ctx context.Context, loc model.VehicleLocation) error {
	return s.repo.Insert(ctx, loc)
}

func (s *vehicleLocationService) GetLastLocation(ctx context.Context, vehicleID string) (*model.VehicleLocation, error) {
	return s.repo.GetLastLocation(ctx, vehicleID)
}

func (s *vehicleLocationService) GetHistory(ctx context.Context, vehicleID string, start, end time.Time) ([]model.VehicleLocation, error) {
	return s.repo.GetHistory(ctx, vehicleID, start, end)
}
