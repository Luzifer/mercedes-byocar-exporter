package mercedes

import (
	"fmt"

	"github.com/pkg/errors"
)

type (
	ElectricStatus struct {
		// Displayed state of charge for the HV battery	0..100 %
		StateOfCharge TimedInt `apiField:"soc"`
		// Electric range	0..2046 km
		ElectricRange TimedInt `apiField:"rangeelectric"`
	}
)

func (a APIClient) GetElectricStatus(vehicleID string) (ElectricStatus, error) {
	var (
		path = fmt.Sprintf("/vehicles/%s/containers/electricvehicle", vehicleID)
		out  ElectricStatus
	)

	if err := a.request(path, &out); err != nil {
		return out, errors.Wrap(err, "getting electric status")
	}

	return out, nil
}
