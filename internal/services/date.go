package services

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

type Date struct {
	TimeZone      string
	ISOStringDate string
	HourMinSec    string //hh:mm:ss
}

const TimeLayout string = "2006-01-02 15:04:05.999999999 -0700 MST"
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

	ISOTime, err := time.Parse(TimeLayout, d.ISOStringDate)
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
// "2006-01-02 15:04:05.999999999 -0700 MST"
func (d *Date) hourMinSecStr() string {
	currentTime := time.Now()

	year := currentTime.Format("2006")
	month := currentTime.Format("01")
	day := currentTime.Format("02")
	milliseconds := currentTime.Format(".000000")
	offset := currentTime.Format("-0700")
	zoneAbbr := currentTime.Format("MST")  // time zone abbreviation
	date := year + "-" + month + "-" + day // 2006-01-02
	time := d.HourMinSec + milliseconds    // 15:04:05.999999999
	hourMinSecStr := date + " " + time + " " + offset + " " + zoneAbbr

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

// Returns the layout in the format
// "2006-01-02 15:04:05.999999999 -0700 MST"
// by extracting offset and time zone abbreviation
// from the ISOString input of the same format
func (d *Date) ISOTimeLayout() (string, error) {
	dateStr := "2006-01-02 15:04:05.999999999"

	offset, err := d.extractOffset(d.ISOStringDate)
	if err != nil {
		fmt.Println("Error Extracting offset:", err)
		return "", err
	}

	timeZoneAbbreviation, err := d.extractTimeZoneAbbreviation(d.ISOStringDate)
	if err != nil {
		fmt.Println("Error Extracting time zone abbreviation:", err)
		return "", err
	}

	ISOString := dateStr + " " + offset + " " + timeZoneAbbreviation

	return ISOString, nil
}

// Returns the layout in the format
// "2006-Jan-02 15:04:05 -0700" by extracting offset
// from the ISOString input of the same format
func (d *Date) TimeLayout() (string, error) {
	hourMinSec := "15:04:05"
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
	offset, err := d.extractOffset(d.ISOStringDate)
	if err != nil {
		fmt.Println("Error Extracting offset:", err)
		return "", err
	}

	date := currentTimeYear + "-" + currentTimeMonth + "-" + currentTimeDay
	layoutStr := date + " " + hourMinSec + " " + offset

	return layoutStr, nil
}

// Returns the layout in the format
// "2006-01-02T15:04:05.999999999+03:00" by extracting offset
// from the RFC3339Nano input of the same format
func (d *Date) RFC3339NanoLayout() (string, error) {
	dateStr := "2006-01-02T15:04:05.999999999"

	offset, err := d.extractOffsetFromRFC3339Nano(d.ISOStringDate)
	if err != nil {
		fmt.Println("Error Extracting offset:", err)
		return "", err
	}
	RFC3339NanoString := dateStr + offset

	return RFC3339NanoString, nil
}

// extractOffset extracts the
// offset from  a given time string of format
// "2006-01-02 15:04:05.999999999 -0700 MST" or
// "2006-Jan-02 15:04:05 -0700" using string matching.
func (d *Date) extractOffset(timeStr string) (string, error) {
	re := regexp.MustCompile(`[-+]\d{4}`)
	offset := re.FindString(timeStr)
	if offset == "" {
		return "", errors.New("offset not found in time string")
	}
	return offset, nil
}

// extractTimeZoneAbbreviation extracts the
// time zone abbreviation from a given time string
// of format "2006-01-02 15:04:05.999999999 -0700 MST"
// using string matching.
func (d *Date) extractTimeZoneAbbreviation(timeStr string) (string, error) {
	re := regexp.MustCompile(`[A-Z]{3,4}$`)
	abbr := re.FindString(timeStr)
	if abbr == "" {
		return "", errors.New("time zone abbreviation not found in time string")
	}
	return abbr, nil
}

// extractOffsetFromRFC3339Nano extracts
// the offset from a given RFC3339Nano formatted time
// string "2006-01-02T15:04:05.999999999+03:00"	using string matching.
func (d *Date) extractOffsetFromRFC3339Nano(timeStr string) (string, error) {
	re := regexp.MustCompile(`[-+]\d{2}:\d{2}`)
	offset := re.FindString(timeStr)
	if offset == "" {
		return "", errors.New("offset not found in time string")
	}
	return offset, nil
}
