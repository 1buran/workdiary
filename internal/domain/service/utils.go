package service

import (
	"time"
)

// Get weekday number in new style: Monday (1st day) is the begining of week,
// Tuesday - 2d day and so on: Mon = 1, Tue = 2 ... Sun = 7
func GetWeekDayNumber(t time.Time) int {
	v := int(t.Weekday())
	if v == 0 {
		return 7
	}
	return v
}
