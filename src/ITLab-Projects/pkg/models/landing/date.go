package landing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Time struct {
	time.Time	`json:",inline" bson:",inline"`
}

func (t Time) MarshalJSON() ([]byte, error) {
	date := fmt.Sprintf(
		"%v/%v/%v",
		t.Day(),
		int(t.Month()),
		t.Year(),
	)

	return json.Marshal(date)
}

func (t *Time) UnmarshalJSON(data []byte) error {
	date := bytes.Split(bytes.Trim(data, `"`), []byte{'/'})

	if len(date) != 3 {
		return fmt.Errorf("Not valid date")
	}

	dayString := string(date[0])
	monthString := string(date[1])
	yearString := string(date[2])

	day, err := strconv.ParseInt(dayString, 10, 64)
	if err != nil {
		return fmt.Errorf("Not valid date")
	}

	month, err := strconv.ParseInt(monthString, 10, 64)
	if err != nil {
		return fmt.Errorf("Not valid date")
	}

	year, err := strconv.ParseInt(yearString, 10, 64)
	if err != nil {
		return fmt.Errorf("Not valid date")
	}

	t.Time = time.Date(
		int(year),
		time.Month(month),
		int(day),
		0,
		0,
		0,
		0,
		time.Local,
	)

	return nil
}