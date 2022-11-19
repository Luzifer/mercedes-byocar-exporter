package mercedes

import (
	"fmt"

	"github.com/pkg/errors"
)

type (
	PayAsYouDriveInsurance struct {
		Odometer TimedInt `apiField:"odo"`
	}
)

func (a APIClient) GetPayAsYouDriveInsurance(vehicleID string) (PayAsYouDriveInsurance, error) {
	var (
		path = fmt.Sprintf("/vehicles/%s/containers/payasyoudrive", vehicleID)
		out  PayAsYouDriveInsurance
	)

	if err := a.request(path, &out); err != nil {
		return out, errors.Wrap(err, "getting pay-as-you-drive response")
	}

	return out, nil
}
