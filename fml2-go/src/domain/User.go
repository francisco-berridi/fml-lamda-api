package domain

import (
	"enums/Unit"
	"time"
	"utils/conversions"
)

type User struct {
	Id             string `json:"id"`
	Username       string `json:"username"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	UserAttributes `json:"userAttributes,omitempty"`
	SyncData
	Settings
}

type SyncData struct {
	LastSynchronization *time.Time `json:"lastSync,omitempty"`
}

type Settings struct {
	HeightUnits Unit.Enum
	WeightUnits Unit.Enum
}

func (u User) SetHeight(value int64, unit Unit.Enum) {

	if unit == Unit.Cm {
		u.HeightCM = value
	} else if unit == Unit.Ft {
		u.HeightCM = int64(conversions.FeetToCentimeters(float64(value)))
	}

	u.HeightUnits = unit
}
func (u User) SetWeight(value float64, unit Unit.Enum) {
	if unit == Unit.Kg {
		u.WeightKG = value
	} else if unit == Unit.Lb {
		u.WeightKG = conversions.LbToKg(value)
	}

	u.WeightUnits = unit
}
