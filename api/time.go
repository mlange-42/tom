package api

import (
	"fmt"
	"strings"
	"time"
)

const timeLayout = "2006-01-02T15:04"
const dateLayout = "2006-01-02"

type Time struct {
	time.Time
}

type Date struct {
	time.Time
}

func (ct *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	ct.Time, err = time.Parse(timeLayout, s)
	return
}

func (ct *Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(timeLayout))), nil
}

func (ct *Date) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	ct.Time, err = time.Parse(dateLayout, s)
	return
}

func (ct *Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(dateLayout))), nil
}
