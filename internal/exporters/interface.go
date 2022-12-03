package exporters

import "github.com/Luzifer/mercedes-byocar-exporter/internal/mercedes"

type (
	Exporter interface {
		SetElectricStatus(vehicleID string, es mercedes.ElectricStatus)
		SetFuelStatus(vehicleID string, fs mercedes.FuelStatus)
		SetLockStatus(vehicleID string, ls mercedes.LockStatus)
		SetPayAsYouGo(vehicleID string, p mercedes.PayAsYouDriveInsurance)
		SetVehicleStatus(vehicleID string, vs mercedes.VehicleStatus)
	}

	Set []Exporter
)

var _ Exporter = Set{}

func (s Set) SetElectricStatus(vehicleID string, es mercedes.ElectricStatus) {
	for _, e := range s {
		e.SetElectricStatus(vehicleID, es)
	}
}

func (s Set) SetFuelStatus(vehicleID string, fs mercedes.FuelStatus) {
	for _, e := range s {
		e.SetFuelStatus(vehicleID, fs)
	}
}

func (s Set) SetLockStatus(vehicleID string, ls mercedes.LockStatus) {
	for _, e := range s {
		e.SetLockStatus(vehicleID, ls)
	}
}

func (s Set) SetPayAsYouGo(vehicleID string, p mercedes.PayAsYouDriveInsurance) {
	for _, e := range s {
		e.SetPayAsYouGo(vehicleID, p)
	}
}

func (s Set) SetVehicleStatus(vehicleID string, vs mercedes.VehicleStatus) {
	for _, e := range s {
		e.SetVehicleStatus(vehicleID, vs)
	}
}
