package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	labelVehicleID = "vehicle_id"
	labelDoor      = "door"
	labelLight     = "light"
	labelWindow    = "window"

	metricsNamespace = "mercedes_byocar"

	subsystemElectricStatus = `electric_status`
	subsystemFuelStatus     = "fuel_status"
	subsystemLockStatus     = "lock_status"
	subsystemPayAsYouDrive  = "pay_as_you_drive"
	subsystemVehicleStatus  = "vehicle_status"
)

var (
	electricSOC   *prometheus.GaugeVec
	electricRange *prometheus.GaugeVec

	fuelRangeLiquidVec   *prometheus.GaugeVec
	fuelTanklevelPercent *prometheus.GaugeVec

	lockDeckLidUnlocked *prometheus.GaugeVec
	lockVehicleStatus   *prometheus.GaugeVec
	lockGasLidUnlocked  *prometheus.GaugeVec
	lockHeading         *prometheus.GaugeVec

	paydOdometer *prometheus.GaugeVec

	vehicleDeckLidOpen   *prometheus.GaugeVec
	vehicleDoorOpen      *prometheus.GaugeVec
	vehicleInteriorLight *prometheus.GaugeVec
	vehicleLightSwitch   *prometheus.GaugeVec
	vehicleReadingLampOn *prometheus.GaugeVec
	vehicleRoofTopStatus *prometheus.GaugeVec
	vehicleSunRoofStatus *prometheus.GaugeVec
	vehicleWindowStatus  *prometheus.GaugeVec
)

func init() {
	initElectricStatus()
	initFuelStatus()
	initLockStatus()
	initPAYD()
	initVehicleStatus()
}

func initElectricStatus() {
	electricRange = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemElectricStatus,
		Name:      "electric_range",
		Help:      "Electric range - 0..2046 km",
	}, []string{labelVehicleID})

	electricSOC = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemElectricStatus,
		Name:      "state_of_charge",
		Help:      "Displayed state of charge for the HV battery - 0..100 %",
	}, []string{labelVehicleID})
}

func initFuelStatus() {
	fuelRangeLiquidVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemFuelStatus,
		Name:      "range_liquid",
		Help:      "Liquid fuel tank range - 0..2046 km",
	}, []string{labelVehicleID})

	fuelTanklevelPercent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemFuelStatus,
		Name:      "tanklevel_percent",
		Help:      "Liquid fuel tank level - 0..100 %",
	}, []string{labelVehicleID})
}

func initLockStatus() {
	lockDeckLidUnlocked = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemLockStatus,
		Name:      "deck_lid_unlocked",
		Help:      "Lock status of the deck lid - 1 = unlocked",
	}, []string{labelVehicleID})

	lockVehicleStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemLockStatus,
		Name:      "vehicle_status",
		Help:      "Vehicle lock status - 0 = unlocked, 1 = internal locked, 2 = external locked, 3 = selective unlocked",
	}, []string{labelVehicleID})

	lockGasLidUnlocked = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemLockStatus,
		Name:      "gas_lid_unlocked",
		Help:      "Status of gas tank door lock - 1 = unlocked",
	}, []string{labelVehicleID})

	lockHeading = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemLockStatus,
		Name:      "heading",
		Help:      "Vehicle heading position - 0..359.9 degrees",
	}, []string{labelVehicleID})
}

func initPAYD() {
	paydOdometer = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemPayAsYouDrive,
		Name:      "odometer",
		Help:      "Odometer - 0..999999 km",
	}, []string{labelVehicleID})
}

func initVehicleStatus() {
	vehicleDeckLidOpen = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemVehicleStatus,
		Name:      "deck_lid_open",
		Help:      "Deck lid latch status opened/closed state - 1 = open",
	}, []string{labelVehicleID})

	vehicleDoorOpen = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemVehicleStatus,
		Name:      "door_open",
		Help:      "Status of respective door - 1 = open",
	}, []string{labelVehicleID, labelDoor})

	vehicleInteriorLight = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemVehicleStatus,
		Name:      "interior_light_on",
		Help:      "Status of respective interior light - 1 = on",
	}, []string{labelVehicleID, labelLight})

	vehicleLightSwitch = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemVehicleStatus,
		Name:      "light_switch_position",
		Help:      "Rotary light switch position - 0 = auto, 1 = headlights, 2 = sidelight left, 3 = sidelight right, 4 = parking light",
	}, []string{labelVehicleID})

	vehicleReadingLampOn = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemVehicleStatus,
		Name:      "reading_lamp_on",
		Help:      "Status of respective reading lamp - 1 = on",
	}, []string{labelVehicleID, labelLight})

	vehicleRoofTopStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemVehicleStatus,
		Name:      "roof_top_status",
		Help:      "Status of the convertible top - 0 = unlocked, 1 = open and locked, 2 = closed and locked",
	}, []string{labelVehicleID})

	vehicleSunRoofStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemVehicleStatus,
		Name:      "sun_roof_status",
		Help:      "Status of the sunroof - 0 = Tilt/slide sunroof is closed, 1 = Tilt/slide sunroof is complete open, 2 = Lifting roof is open, 3 = Tilt/slide sunroof is running, 4 = Tilt/slide sunroof in anti-booming position, 5 = Sliding roof in intermediate position, 6 = Lifting roof in intermediate position",
	}, []string{labelVehicleID})

	vehicleWindowStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: subsystemVehicleStatus,
		Name:      "window_status",
		Help:      "Status of respective window - 0 = window in intermediate position, 1 = window completely opened, 2 = window completely closed, 3 = window airing position, 4 = window intermediate airing position, 5 = window currently running",
	}, []string{labelVehicleID, labelWindow})
}
