package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"transjakarta/service"

	"github.com/gofiber/fiber/v2"
)

type VehicleHandler struct {
	service service.VehicleLocationService
}

func NewVehicleHandler(svc service.VehicleLocationService) *VehicleHandler {
	return &VehicleHandler{service: svc}
}

func (h *VehicleHandler) GetLastLocation(c *fiber.Ctx) error {
	vehicleID := c.Params("vehicle_id")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	loc, err := h.service.GetLastLocation(ctx, vehicleID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "data not found"})
	}
	return c.JSON(loc)
}

func (h *VehicleHandler) GetHistory(c *fiber.Ctx) error {
	vehicleID := c.Params("vehicle_id")
	startStr := c.Query("start")
	endStr := c.Query("end")

	startUnix, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid start timestamp"})
	}
	endUnix, err := strconv.ParseInt(endStr, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid end timestamp"})
	}

	start := time.Unix(startUnix, 0)
	end := time.Unix(endUnix, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	history, err := h.service.GetHistory(ctx, vehicleID, start, end)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch history"})
	}
	return c.JSON(history)
}
