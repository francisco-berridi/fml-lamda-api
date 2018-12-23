package domain

import (
	"time"
	"fmt"
)

type FMLTime struct {
	time.Time
}

func (f FMLTime) MarshalJSON()([]byte, error) {

	fmt.Println("is zero", f.IsZero())

	if f.IsZero() {
		return []byte(""), nil
	}

	bytes, err := f.Time.MarshalJSON()
	fmt.Println(string(bytes), err)

	return f.Time.MarshalJSON()
}