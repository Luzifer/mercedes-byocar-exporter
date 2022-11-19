package mercedes

import (
	"fmt"

	"github.com/pkg/errors"
)

type (
	VehicleStatus struct {
		// Deck lid latch status opened/closed state	false: closed / true: open
		DeckLidOpen TimedBool `apiField:"decklidstatus"`
		// Status of the front left door	false: closed / true: open
		DoorFrontLeftOpen TimedBool `apiField:"doorstatusfrontleft"`
		// Status of the front right door	false: closed / true: open
		DoorFrontRightOpen TimedBool `apiField:"doorstatusfrontright"`
		// Status of the rear left door	false: closed / true: open
		DoorRearLeftOpen TimedBool `apiField:"doorstatusrearleft"`
		// Status of the rear right door	false: closed / true: open
		DoorRearRightOpen TimedBool `apiField:"doorstatusrearright"`
		// Front light inside	false: off / true: on
		InteriorLightsFrontOn TimedBool `apiField:"interiorLightsFront"`
		// Rear light inside	false: off / true: on
		InteriorLightsRearOn TimedBool `apiField:"interiorLightsRear"`
		// Rotary light switch position
		LightSwitchPosition TimedEnum `apiField:"lightswitchposition" values:"auto,headlights,sidelight left,sidelight right,parking light"`
		// Front left reading light inside	false: off / true: on
		ReadingLampFrontLeftOn TimedBool `apiField:"readingLampFrontLeft"`
		// Front right reading light inside	false: off / true: on
		ReadingLampFrontRightOn TimedBool `apiField:"readingLampFrontRight"`
		// Status of the convertible top opened/closed
		RoofTopStatus TimedEnum `apiField:"rooftopstatus" values:"unlocked,open and locked,closed and locked"`
		// Status of the sunroof
		SunRoofStatus TimedEnum `apiField:"sunroofstatus" values:"Tilt/slide sunroof is closed,Tilt/slide sunroof is complete open,Lifting roof is open,Tilt/slide sunroof is running,Tilt/slide sunroof in anti-booming position,Sliding roof in intermediate position,Lifting roof in intermediate position"`
		// Status of the front left window
		WindowStatusFrontLeft TimedEnum `apiField:"windowstatusfrontleft" values:"window in intermediate position,window completely opened,window completely closed,window airing position,window intermediate airing position,window currently running"`
		// Status of the front right window
		WindowStatusFrontRight TimedEnum `apiField:"windowstatusfrontright" values:"window in intermediate position,window completely opened,window completely closed,window airing position,window intermediate airing position,window currently running"`
		// Status of the rear left window
		WindowStatusRearLeft TimedEnum `apiField:"windowstatusrearleft" values:"window in intermediate position,window completely opened,window completely closed,window airing position,window intermediate airing position,window currently running"`
		// Status of the rear right window
		WindowStatusRearRight TimedEnum `apiField:"windowstatusrearright" values:"window in intermediate position,window completely opened,window completely closed,window airing position,window intermediate airing position,window currently running"`
	}
)

func (a APIClient) GetVehicleStatus(vehicleID string) (VehicleStatus, error) {
	var (
		path = fmt.Sprintf("/vehicles/%s/containers/vehiclestatus", vehicleID)
		out  VehicleStatus
	)

	if err := a.request(path, &out); err != nil {
		return out, errors.Wrap(err, "getting vehicle status")
	}

	return out, nil
}
