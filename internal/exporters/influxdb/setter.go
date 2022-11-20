package influxdb

import (
	"strings"

	"github.com/Luzifer/mercedes-byocar-exporter/internal/mercedes"
)

const (
	labelVehicleID = "vehicle_id"
	labelDoor      = "door"
	labelLight     = "light"
	labelWindow    = "window"

	subsystemFuelStatus    = "fuel_status"
	subsystemLockStatus    = "lock_status"
	subsystemPayAsYouDrive = "pay_as_you_drive"
	subsystemVehicleStatus = "vehicle_status"
)

func (e *Exporter) SetFuelStatus(vehicleID string, fs mercedes.FuelStatus) {
	e.submitValue(fs.RangeLiquid, mn(subsystemFuelStatus, "range_liquid"), labelVehicleID, vehicleID)
	e.submitValue(fs.TanklevelPercent, mn(subsystemFuelStatus, "tanklevel_percent"), labelVehicleID, vehicleID)
}

func (e *Exporter) SetLockStatus(vehicleID string, ls mercedes.LockStatus) {
	e.submitValue(ls.DeckLidUnlocked, mn(subsystemLockStatus, "deck_lid_unlocked"), labelVehicleID, vehicleID)
	e.submitValue(ls.VehicleStatus, mn(subsystemLockStatus, "vehicle_status"), labelVehicleID, vehicleID)
	e.submitValue(ls.GasLidUnlocked, mn(subsystemLockStatus, "gas_lid_unlocked"), labelVehicleID, vehicleID)
	e.submitValue(ls.Heading, mn(subsystemLockStatus, "heading"), labelVehicleID, vehicleID)
}

func (e *Exporter) SetPayAsYouGo(vehicleID string, p mercedes.PayAsYouDriveInsurance) {
	e.submitValue(p.Odometer, mn(subsystemPayAsYouDrive, "odometer"), labelVehicleID, vehicleID)
}

func (e *Exporter) SetVehicleStatus(vehicleID string, vs mercedes.VehicleStatus) {
	e.submitValue(vs.DeckLidOpen, mn(subsystemVehicleStatus, "deck_lid_open"), labelVehicleID, vehicleID)

	e.submitValue(vs.DoorFrontLeftOpen, mn(subsystemVehicleStatus, "door_open"), labelVehicleID, vehicleID, labelDoor, "front_left")
	e.submitValue(vs.DoorFrontRightOpen, mn(subsystemVehicleStatus, "door_open"), labelVehicleID, vehicleID, labelDoor, "front_right")
	e.submitValue(vs.DoorRearLeftOpen, mn(subsystemVehicleStatus, "door_open"), labelVehicleID, vehicleID, labelDoor, "rear_left")
	e.submitValue(vs.DoorRearRightOpen, mn(subsystemVehicleStatus, "door_open"), labelVehicleID, vehicleID, labelDoor, "rear_right")

	e.submitValue(vs.InteriorLightsFrontOn, mn(subsystemVehicleStatus, "interior_light_on"), labelVehicleID, vehicleID, labelLight, "front")
	e.submitValue(vs.InteriorLightsRearOn, mn(subsystemVehicleStatus, "interior_light_on"), labelVehicleID, vehicleID, labelLight, "rear")

	e.submitValue(vs.LightSwitchPosition, mn(subsystemVehicleStatus, "light_switch_position"), labelVehicleID, vehicleID)

	e.submitValue(vs.ReadingLampFrontLeftOn, mn(subsystemVehicleStatus, "reading_lamp_on"), labelVehicleID, vehicleID, labelLight, "front_left")
	e.submitValue(vs.ReadingLampFrontRightOn, mn(subsystemVehicleStatus, "reading_lamp_on"), labelVehicleID, vehicleID, labelLight, "front_right")

	e.submitValue(vs.RoofTopStatus, mn(subsystemVehicleStatus, "roof_top_status"), labelVehicleID, vehicleID)
	e.submitValue(vs.SunRoofStatus, mn(subsystemVehicleStatus, "sun_roof_status"), labelVehicleID, vehicleID)

	e.submitValue(vs.WindowStatusFrontLeft, mn(subsystemVehicleStatus, "window_status"), labelVehicleID, vehicleID, labelWindow, "front_left")
	e.submitValue(vs.WindowStatusFrontRight, mn(subsystemVehicleStatus, "window_status"), labelVehicleID, vehicleID, labelWindow, "front_right")
	e.submitValue(vs.WindowStatusRearLeft, mn(subsystemVehicleStatus, "window_status"), labelVehicleID, vehicleID, labelWindow, "rear_left")
	e.submitValue(vs.WindowStatusRearRight, mn(subsystemVehicleStatus, "window_status"), labelVehicleID, vehicleID, labelWindow, "rear_right")
}

func (e *Exporter) submitValue(value mercedes.MetricValue, metric_name string, tvs ...string) {
	if !value.IsValid() {
		return
	}

	v := map[string]any{"value": value.ToFloat()}
	e.RecordPoint(metric_name, tags(tvs...), v, value.Time())
}

func mn(parts ...string) string {
	return strings.Join(parts, "_")
}

func tags(kvs ...string) map[string]string {
	out := make(map[string]string)

	if len(kvs)%2 != 0 {
		panic("invalid tags given")
	}

	for i := 0; i < len(kvs); i += 2 {
		out[kvs[i]] = kvs[i+1]
	}

	return out
}
