package main

import (
	"github.com/sirupsen/logrus"

	"github.com/Luzifer/mercedes-byocar-exporter/internal/mercedes"
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
	enabledExporters.SetPayAsYouGo(vehicleID, s1)

	s2, err := mc.GetFuelStatus(cfg.VehicleID[0])
	if err != nil {
		logger.WithError(err).Error("fetching fuel-status data")
		return
	}
	enabledExporters.SetFuelStatus(vehicleID, s2)

	s3, err := mc.GetVehicleStatus(cfg.VehicleID[0])
	if err != nil {
		logger.WithError(err).Error("fetching vehicle-status data")
		return
	}
	enabledExporters.SetVehicleStatus(vehicleID, s3)

	s4, err := mc.GetLockStatus(cfg.VehicleID[0])
	if err != nil {
		logger.WithError(err).Error("fetching lock-status data")
		return
	}
	enabledExporters.SetLockStatus(vehicleID, s4)

	s5, err := mc.GetElectricStatus(cfg.VehicleID[0])
	if err != nil {
		logger.WithError(err).Error("fetching electric-status data")
		return
	}
	enabledExporters.SetElectricStatus(vehicleID, s5)

	logger.Info("data updated successfully")
}
