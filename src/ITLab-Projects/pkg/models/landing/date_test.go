package landing_test

import (
	"encoding/json"
	"testing"
	"time"


	"github.com/ITLab-Projects/pkg/models/landing"
)

func TestFunc_marshallJSON(t *testing.T) {
	var date landing.Time

	date.Time = time.Date(
		2007,
		time.September,
		3,
		0,
		0,
		0,
		0,
		time.Local,
	)

	data, err := json.Marshal(date)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if eq := string(data) == `"03/09/2007"`; !eq {
		t.Log("Assert error")
		t.Log(string(data))
		t.FailNow()
	}

	var newData landing.Time

	json.Unmarshal(data, &newData)

	t.Log(newData.String())
}

func TestFunc_MarshallJSONCustomStruct(t *testing.T) {
	type Custom struct {
		Date landing.Time `json:"date"`
	}

	c := &Custom{
		Date: landing.Time{
			Time: time.Date(
				2007,
				time.September,
				3,
				0,
				0,
				0,
				0,
				time.Local,
			),
		},
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if eq := string(data) == `{"date":"03/09/2007"}`; !eq {
		t.Log("Assert error")
		t.Log(string(data))
		t.FailNow()
	}
}