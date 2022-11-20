package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/Luzifer/mercedes-byocar-exporter/internal/exporters"
	"github.com/Luzifer/mercedes-byocar-exporter/internal/mercedes"
)

type (
	exporter struct{}
)

var (
	Exporter exporter
	_        exporters.Exporter = exporter{}
)

func (exporter) SetFuelStatus(vehicleID string, fs mercedes.FuelStatus) {
	setGaugeVecValue(fs.RangeLiquid, fuelRangeLiquidVec, labelVehicleID, vehicleID)
	setGaugeVecValue(fs.TanklevelPercent, fuelTanklevelPercent, labelVehicleID, vehicleID)
}

func (exporter) SetLockStatus(vehicleID string, ls mercedes.LockStatus) {
	setGaugeVecValue(ls.DeckLidUnlocked, lockDeckLidUnlocked, labelVehicleID, vehicleID)
	setGaugeVecValue(ls.VehicleStatus, lockVehicleStatus, labelVehicleID, vehicleID)
	setGaugeVecValue(ls.GasLidUnlocked, lockGasLidUnlocked, labelVehicleID, vehicleID)
	setGaugeVecValue(ls.Heading, lockHeading, labelVehicleID, vehicleID)
}

func (exporter) SetPayAsYouGo(vehicleID string, p mercedes.PayAsYouDriveInsurance) {
	setGaugeVecValue(p.Odometer, paydOdometer, labelVehicleID, vehicleID)
}

func (exporter) SetVehicleStatus(vehicleID string, vs mercedes.VehicleStatus) {
	setGaugeVecValue(vs.DeckLidOpen, vehicleDeckLidOpen, labelVehicleID, vehicleID)

	setGaugeVecValue(vs.DoorFrontLeftOpen, vehicleDoorOpen, labelVehicleID, vehicleID, labelDoor, "front_left")
	setGaugeVecValue(vs.DoorFrontRightOpen, vehicleDoorOpen, labelVehicleID, vehicleID, labelDoor, "front_right")
	setGaugeVecValue(vs.DoorRearLeftOpen, vehicleDoorOpen, labelVehicleID, vehicleID, labelDoor, "rear_left")
	setGaugeVecValue(vs.DoorRearRightOpen, vehicleDoorOpen, labelVehicleID, vehicleID, labelDoor, "rear_right")

	setGaugeVecValue(vs.InteriorLightsFrontOn, vehicleInteriorLight, labelVehicleID, vehicleID, labelLight, "front")
	setGaugeVecValue(vs.InteriorLightsRearOn, vehicleInteriorLight, labelVehicleID, vehicleID, labelLight, "rear")

	setGaugeVecValue(vs.LightSwitchPosition, vehicleLightSwitch, labelVehicleID, vehicleID)

	setGaugeVecValue(vs.ReadingLampFrontLeftOn, vehicleReadingLampOn, labelVehicleID, vehicleID, labelLight, "front_left")
	setGaugeVecValue(vs.ReadingLampFrontRightOn, vehicleReadingLampOn, labelVehicleID, vehicleID, labelLight, "front_right")

	setGaugeVecValue(vs.RoofTopStatus, vehicleRoofTopStatus, labelVehicleID, vehicleID)
	setGaugeVecValue(vs.SunRoofStatus, vehicleSunRoofStatus, labelVehicleID, vehicleID)

	setGaugeVecValue(vs.WindowStatusFrontLeft, vehicleWindowStatus, labelVehicleID, vehicleID, labelWindow, "front_left")
	setGaugeVecValue(vs.WindowStatusFrontRight, vehicleWindowStatus, labelVehicleID, vehicleID, labelWindow, "front_right")
	setGaugeVecValue(vs.WindowStatusRearLeft, vehicleWindowStatus, labelVehicleID, vehicleID, labelWindow, "rear_left")
	setGaugeVecValue(vs.WindowStatusRearRight, vehicleWindowStatus, labelVehicleID, vehicleID, labelWindow, "rear_right")
}

func setGaugeVecValue(value mercedes.MetricValue, vec *prometheus.GaugeVec, lvs ...string) {
	if !value.IsValid() {
		return
	}

	vec.With(labels(lvs...)).Set(value.ToFloat())
}

func boolToValue(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func labels(kvs ...string) prometheus.Labels {
	out := make(prometheus.Labels)

	if len(kvs)%2 != 0 {
		panic("invalid labels given")
	}

	for i := 0; i < len(kvs); i += 2 {
		out[kvs[i]] = kvs[i+1]
	}

	return out
}
