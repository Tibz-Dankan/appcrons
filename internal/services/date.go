package services

import (
	"fmt"
	"time"
)

type Date struct {
	TimeZone      string
	ISOStringDate string
	HourMinSec    string //hh:mm:ss
}

const ISOStringLayout string = "2006-01-02 15:04:05.999999999 -0700 MST"
const TimeLayout string = "2006-Jan-02 15:04:05"

// Returns the current  in the
// provided time zone
func (d *Date) CurrentTime() (time.Time, error) {
	currentTime := time.Now()

	locTimeZone, err := time.LoadLocation(d.TimeZone)
	if err != nil {
		fmt.Println("Error loading location:", err)
		return time.Now(), err
	}

	return currentTime.In(locTimeZone), nil
}

// Returns the ISOStringDate as time in the
// provided time zone
func (d *Date) ISOTime() (time.Time, error) {

	locTimeZone, err := time.LoadLocation(d.TimeZone)
	if err != nil {
		fmt.Println("Error loading location:", err)
		return time.Now(), err
	}

	ISOTime, err := time.Parse(ISOStringLayout, d.ISOStringDate)
	if err != nil {
		return time.Now(), err
	}

	return ISOTime.In(locTimeZone), nil

}

// Return the date string in the format
// "2006-Jan-02 15:04:05"
func (d *Date) hourMinSecStr() string {
	currentTime := time.Now()

	var currentTimeDay string
	day := currentTime.Day()
	if day < 10 {
		currentTimeDay = "0" + fmt.Sprint(day)
	}
	currentTimeYear := fmt.Sprint(currentTime.Year())
	currentTimeMonth := fmt.Sprint(currentTime.Month())

	date := currentTimeYear + "-" + currentTimeMonth + "-" + currentTimeDay
	hourMinSecStr := date + " " + d.HourMinSec

	return hourMinSecStr
}

// Returns the HourMinSec as time in the
// provided time zone
func (d *Date) HourMinSecTime() (time.Time, error) {

	locTimeZone, err := time.LoadLocation(d.TimeZone)
	if err != nil {
		fmt.Println("Error loading location:", err)
		return time.Now(), err
	}

	HourMinSecTime, err := time.Parse(TimeLayout, d.hourMinSecStr())
	if err != nil {
		return time.Now(), err
	}

	return HourMinSecTime.In(locTimeZone), nil
}
