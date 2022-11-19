package mercedes

import (
	"fmt"

	"github.com/pkg/errors"
)

type (
	FuelStatus struct {
		RangeLiquid      TimedInt `apiField:"rangeliquid"`      // Liquid fuel tank range	0..2046 km
		TanklevelPercent TimedInt `apiField:"tanklevelpercent"` // Liquid fuel tank level	0…100 %
	}
)

func (a APIClient) GetFuelStatus(vehicleID string) (FuelStatus, error) {
	var (
		path = fmt.Sprintf("/vehicles/%s/containers/fuelstatus", vehicleID)
		out  FuelStatus
	)

	if err := a.request(path, &out); err != nil {
		return out, errors.Wrap(err, "getting fuel status")
	}

	return out, nil
}
