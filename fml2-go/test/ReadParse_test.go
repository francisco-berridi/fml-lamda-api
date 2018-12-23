package test

import (
	"fmt"
	"services/parsehttp"
	"testing"
)

func TestReadParseWorkoutSchedules(t *testing.T) {

	resultJson, err := parsehttp.GetWorkoutSchedules()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(resultJson)

}
