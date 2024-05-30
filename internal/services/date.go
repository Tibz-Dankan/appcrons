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
const TimeLayout string = "2006-Jan-02 15:04:05 -0700"
const RFC3339NanoLayout string = "2006-01-02T15:04:05.999999999+03:00"

// Returns the current in the provided time zone
// with seconds set to 00
func (d *Date) CurrentTime() (time.Time, error) {
	currentTime := time.Now()

	locTimeZone, err := time.LoadLocation(d.TimeZone)
	if err != nil {
		fmt.Println("Error loading location:", err)
		return time.Now(), err
	}

	// Setting seconds part to "00"
	currentTime = currentTime.In(locTimeZone)
	currentTime = currentTime.Truncate(time.Minute)

	return currentTime, nil
}

// Returns the ISOStringDate as time in the
// provided time zone with seconds set to 00
func (d *Date) ISOTime() (time.Time, error) {

	locTimeZone, err := time.LoadLocation(d.TimeZone)
	if err != nil {
		fmt.Println("Error loading location:", err)
		return time.Now(), err
	}

	ISOTime, err := time.Parse(ISOStringLayout, d.ISOStringDate)
	if err != nil {
		fmt.Println("Error parsing ISOStringDate:", err)
		return time.Now(), err
	}

	// Setting seconds part to "00"
	ISOTime = ISOTime.In(locTimeZone)
	ISOTime = ISOTime.Truncate(time.Minute)

	return ISOTime, nil
}

// Return the date string in the format
// "2006-Jan-02 15:04:05"
func (d *Date) hourMinSecStr() string {
	currentTime := time.Now()

	var currentTimeDay string
	day := currentTime.Day()
	if day < 10 {
		currentTimeDay = "0" + fmt.Sprint(day)
	} else {
		currentTimeDay = fmt.Sprint(day)
	}
	currentTimeYear := fmt.Sprint(currentTime.Year())
	currentTimeMonth := fmt.Sprint(currentTime.Month())
	offset := currentTime.Format("-0700")

	date := currentTimeYear + "-" + currentTimeMonth + "-" + currentTimeDay
	hourMinSecStr := date + " " + d.HourMinSec + " " + offset

	return hourMinSecStr
}

// Returns the HourMinSec as time in the
// provided time zone with seconds set to 00
func (d *Date) HourMinSecTime() (time.Time, error) {

	locTimeZone, err := time.LoadLocation(d.TimeZone)
	if err != nil {
		fmt.Println("Error loading location:", err)
		return time.Now(), err
	}

	HourMinSecTime, err := time.Parse(TimeLayout, d.hourMinSecStr())
	if err != nil {
		fmt.Println("Error parsing hourMinSecStr:", err)
		return time.Now(), err
	}

	// Setting seconds part to "00"
	HourMinSecTime = HourMinSecTime.In(locTimeZone)
	HourMinSecTime = HourMinSecTime.Truncate(time.Minute)

	return HourMinSecTime.In(locTimeZone), nil
}

// Returns time in RFC3339Nano(includes nanoseconds) format
// of the provided IsoStringDate in the format like
// '2024-05-09T13:42:59.994557+03:00'
func (d *Date) RFC3339Nano() (time.Time, error) {

	parsedTime, err := time.Parse(RFC3339NanoLayout, d.ISOStringDate)
	if err != nil {
		fmt.Println("Error parsing RFC3339NanoInput:", err)
		return time.Now(), err
	}

	return parsedTime, nil
}
