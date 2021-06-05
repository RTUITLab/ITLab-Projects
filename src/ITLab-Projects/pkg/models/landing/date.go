package landing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

)

type Time struct {
	time.Time	`json:",inline"`
}

func (t Time) ToStringDate() string {
	monthInt := int(t.Month())
	var month string
	if monthInt < 10 {
		month = fmt.Sprintf("0%v", monthInt)
	} else {
		month = fmt.Sprint(monthInt)
	}

	dayInt := t.Day()
	var day string
	if dayInt < 10 {
		day = fmt.Sprintf("0%v", dayInt)
	} else {
		day = fmt.Sprint(dayInt)
	}
	date := fmt.Sprintf(
		"%v/%s/%v",
		day,
		month,
		t.Year(),
	)

	return date
}

func (t *Time) FromString(dayString, monthString, yearString string) error {
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

func (t *Time) FromBytes(data []byte) error {
	logrus.Info(string(data))
	date := bytes.Split(bytes.Trim(data, `"`), []byte{'/'})
	if len(date) != 3 {
		return fmt.Errorf("Not valid date")
	}

	dayString := string(date[0])
	monthString := string(date[1])
	yearString := string(date[2])

	return t.FromString(dayString, monthString, yearString)
}

func (t Time) MarshalJSON() ([]byte, error) {
	date := t.ToStringDate()

	return json.Marshal(date)
}

func (t *Time) UnmarshalJSON(data []byte) error {
	return t.FromBytes(data)
}