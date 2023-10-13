package helpers

import (
	"errors"
	"time"
)

func ParseDateOrStringToDate(date interface{}) (time.Time, error) {
	if _, ok := date.(string); ok {
		realDate, error := time.Parse("DD-MM-YYYY", date.(string))
		if error != nil {
			panic(error)
		}
		return realDate, nil
	}
	if _, ok := date.(time.Time); ok {

		return date.(time.Time), nil
	}
	return time.Time{}, errors.New("Invalid date type")
}
func GetEndOfDay(date interface{}) time.Time {
	realDate, err := ParseDateOrStringToDate(date)
	if err != nil {
		panic(err)
	}
	return time.Date(realDate.Year(), realDate.Month(), realDate.Day(), 23, 59, 59, 99, realDate.Location())
}

func GetStartOfDay(date time.Time) time.Time {
	realDate, err := ParseDateOrStringToDate(date)
	if err != nil {
		panic(err)
	}

	return time.Date(realDate.Year(), realDate.Month(), realDate.Day(), 0, 0, 0, 0, realDate.Location())
}
