package domain

import (
	"errors"
	"time"
)

func ParseTime(t string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return time.Time{}, err
	}

	if parsedTime.Year() < 1970 {
		return time.Time{}, errors.New("year cannot be less than 1970")
	}
	return parsedTime, nil
}
