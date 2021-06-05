package landing_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ITLab-Projects/pkg/models/landing"
)

func TestFunc_marshallJSON(t *testing.T) {
	var date landing.Time

	date.Time = time.Now()

	data, _ := json.Marshal(date)

	t.Log(string(data))

	var newData landing.Time

	json.Unmarshal(data, &newData)

	t.Log(newData.String())
}