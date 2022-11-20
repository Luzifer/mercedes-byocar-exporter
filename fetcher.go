package main

import (
	"github.com/sirupsen/logrus"

	"github.com/Luzifer/mercedes-byocar-exporter/internal/mercedes"
	"github.com/Luzifer/mercedes-byocar-exporter/internal/prometheus"
)

func getCronFunc(mc mercedes.Client) func() {
	return func() {
		for i := range cfg.VehicleID {
			runFetcher(mc, cfg.VehicleID[i])
		}
	}
}

func runFetcher(mc mercedes.Client, vehicleID string) {
	logger := logrus.WithField("vehicle_id", vehicleID)
	logger.Info("fetching data")

	s1, err := mc.GetPayAsYouDriveInsurance(cfg.VehicleID[0])
	if err != nil {
		logger.WithError(err).Error("fetching pay-as-you-go data")
		return
	}
	prometheus.SetPayAsYouGo(vehicleID, s1)

	s2, err := mc.GetFuelStatus(cfg.VehicleID[0])
	if err != nil {
		logger.WithError(err).Error("fetching fuel-status data")
		return
	}
	prometheus.SetFuelStatus(vehicleID, s2)

	s3, err := mc.GetVehicleStatus(cfg.VehicleID[0])
	if err != nil {
		logger.WithError(err).Error("fetching vehicle-status data")
		return
	}
	prometheus.SetVehicleStatus(vehicleID, s3)

	s4, err := mc.GetLockStatus(cfg.VehicleID[0])
	if err != nil {
		logger.WithError(err).Error("fetching lock-status data")
		return
	}
	prometheus.SetLockStatus(vehicleID, s4)

	logger.Info("data updated successfully")
}
