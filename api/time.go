package api

import (
	"fmt"
	"strings"
	"time"
)

const TimeLayout = "2006-01-02T15:04"
const DateLayout = "2006-01-02"
const DateLayoutShort = "Jan 2"

type Time struct {
	time.Time
}

type Date struct {
	time.Time
}

func (ct *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	ct.Time, err = time.Parse(TimeLayout, s)
	return
}

func (ct *Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(TimeLayout))), nil
}

func (ct *Date) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	ct.Time, err = time.Parse(DateLayout, s)
	return
}

func (ct *Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(DateLayout))), nil
}
