package main

import (
	"errors"

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
	handleMetricsEntries(logger, "pay-as-you-go", err, func() { enabledExporters.SetPayAsYouGo(vehicleID, s1) })

	s2, err := mc.GetFuelStatus(cfg.VehicleID[0])
	handleMetricsEntries(logger, "fuel-status", err, func() { enabledExporters.SetFuelStatus(vehicleID, s2) })

	s3, err := mc.GetVehicleStatus(cfg.VehicleID[0])
	handleMetricsEntries(logger, "vehicle-status", err, func() { enabledExporters.SetVehicleStatus(vehicleID, s3) })

	s4, err := mc.GetLockStatus(cfg.VehicleID[0])
	handleMetricsEntries(logger, "lock-status", err, func() { enabledExporters.SetLockStatus(vehicleID, s4) })

	s5, err := mc.GetElectricStatus(cfg.VehicleID[0])
	handleMetricsEntries(logger, "electric-status", err, func() { enabledExporters.SetElectricStatus(vehicleID, s5) })

	logger.Info("data updated")
}

func handleMetricsEntries(logger *logrus.Entry, dataType string, err error, submit func()) {
	switch {
	case err == nil:
		submit()

	case errors.Is(err, mercedes.ErrNoDataAvailable):
		logger.Warnf("%s data is not available", dataType)
		return

	default:
		logger.WithError(err).Errorf("fetching %s data", dataType)
		return

	}
}
