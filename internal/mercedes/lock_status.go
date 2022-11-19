package mercedes

import (
	"fmt"

	"github.com/pkg/errors"
)

type (
	LockStatus struct {
		// Lock status of the deck lid	false: locked / true: unlocked
		DeckLidUnlocked TimedBool `apiField:"doorlockstatusdecklid"`
		// Vehicle lock status
		VehicleStatus TimedEnum `apiField:"doorlockstatusvehicle" values:"unlocked,internal locked,external locked,selective unlocked"`
		// Status of gas tank door lock	false: locked / true: unlocked
		GasLidUnlocked TimedBool `apiField:"doorlockstatusgas"`
		// Vehicle heading position	0..359.9 degrees
		Heading TimedFloat `apiField:"positionHeading"`
	}
)

func (a APIClient) GetLockStatus(vehicleID string) (LockStatus, error) {
	var (
		path = fmt.Sprintf("/vehicles/%s/containers/vehiclelockstatus", vehicleID)
		out  LockStatus
	)

	if err := a.request(path, &out); err != nil {
		return out, errors.Wrap(err, "getting lock status")
	}

	return out, nil
}
