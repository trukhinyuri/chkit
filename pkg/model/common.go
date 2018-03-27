package model

import (
	"fmt"
	"time"
)

const (
	CreationTimeFormat = time.RFC822
)

func TimestampFormat(timestamp time.Time) string {
	age := time.Now().Sub(timestamp)
	ageString := age.String()
	const year = 365 * 24
	switch {
	case age.Hours() > year:
		years := uint64(age.Hours()) / year
		ageString = fmt.Sprintf("%dy", years)
	case age.Hours() > 24 && age.Hours() < 365*24:
		days := uint64(age.Hours()) / 24
		ageString = fmt.Sprintf("%dd", days)
	case age.Hours() <= 24 && age.Hours() > 1:
		hours := uint64(age.Hours())
		ageString = fmt.Sprintf("%dh", hours)
	case age.Hours() < 1 && age.Minutes() > 1:
		minutes := uint64(age.Minutes())
		ageString = fmt.Sprintf("%dm", minutes)
	default:
		seconds := uint64(age.Seconds())
		ageString = fmt.Sprintf("%ds", seconds)
	}
	return ageString
}
