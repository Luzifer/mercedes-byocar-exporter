package mercedes

import (
	"fmt"
	"net/http"
	"time"
)

type (
	Client interface {
		GetAuthStartURL(redirectURL string) string
		GetFuelStatus(vehicleID string) (FuelStatus, error)
		GetLockStatus(vehicleID string) (LockStatus, error)
		GetPayAsYouDriveInsurance(vehicleID string) (PayAsYouDriveInsurance, error)
		GetVehicleStatus(vehicleID string) (VehicleStatus, error)
		StoreTokenFromRequest(redirectURL string, r *http.Request) error
	}

	MetricValue interface {
		IsValid() bool
		ToFloat() float64
	}

	TimedBool struct {
		v bool
		t time.Time
	}

	TimedEnum struct {
		v   int64
		def []string
		t   time.Time
	}

	TimedFloat struct {
		v float64
		t time.Time
	}

	TimedInt struct {
		v int64
		t time.Time
	}

	genericAPIResponse []map[string]*metricValue

	metricValue struct {
		Value     string `json:"value"`
		Timestamp int64  `json:"timestamp"`
	}
)

var (
	_ MetricValue = TimedBool{}
	_ MetricValue = TimedEnum{}
	_ MetricValue = TimedFloat{}
	_ MetricValue = TimedInt{}
)

const (
	apiPrefix = "https://api.mercedes-benz.com/vehicledata/v2"

	oAuthEndpointAuth  = "https://ssoalpha.dvb.corpinter.net/v1/auth"
	oAuthEndpointToken = "https://ssoalpha.dvb.corpinter.net/v1/token"

	oAuthScopeOfflineAccess     = "offline_access"
	oAuthScopeOpenID            = "openid"
	oAuthScopePayAsYouDrive     = "mb:vehicle:mbdata:payasyoudrive"
	oAuthScopeVehicleFuelStatus = "mb:vehicle:mbdata:fuelstatus"
	oAuthScopeVehicleLockStatus = "mb:vehicle:mbdata:vehiclelock"
	oAuthScopeVehicleStatus     = "mb:vehicle:mbdata:vehiclestatus"
)

func (g genericAPIResponse) Get(key string) *metricValue {
	for i := range g {
		if g[i][key] != nil {
			return g[i][key]
		}
	}

	return nil
}

// Bool

func (t TimedBool) Bool() bool { return t.v }

func (t TimedBool) IsValid() bool { return !t.t.IsZero() }

func (t TimedBool) String() string { return fmt.Sprintf("%v (%s)", t.v, t.t.Format(time.RFC3339)) }

func (t TimedBool) ToFloat() float64 {
	if t.v {
		return 1
	}
	return 0
}

// Enum

func (t TimedEnum) Idx() int64 { return t.v }

func (t TimedEnum) IsValid() bool { return !t.t.IsZero() }

func (t TimedEnum) String() string {
	s := "n/a"
	if len(t.def) > 0 {
		s = t.def[t.v]
	}
	return fmt.Sprintf("%s (%s)", s, t.t.Format(time.RFC3339))
}

func (t TimedEnum) ToFloat() float64 { return float64(t.v) }

func (t TimedEnum) Value() string {
	if len(t.def) > 0 {
		return t.def[t.v]
	}

	return "n/a"
}

// Float

func (t TimedFloat) Float() float64 { return t.v }

func (t TimedFloat) IsValid() bool { return !t.t.IsZero() }

func (t TimedFloat) String() string { return fmt.Sprintf("%v (%s)", t.v, t.t.Format(time.RFC3339)) }

func (t TimedFloat) ToFloat() float64 { return t.v }

// Int

func (t TimedInt) Int() int64 { return t.v }

func (t TimedInt) IsValid() bool { return !t.t.IsZero() }

func (t TimedInt) String() string { return fmt.Sprintf("%v (%s)", t.v, t.t.Format(time.RFC3339)) }

func (t TimedInt) ToFloat() float64 { return float64(t.v) }
