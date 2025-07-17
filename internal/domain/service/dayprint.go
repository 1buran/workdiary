package service

import (
	"fmt"
	"math"

	"github.com/muesli/termenv"

	"github.com/1buran/workdiary/internal/domain/valueobject"
)

func NewDayPrinter(
	output *termenv.Output,
	paletter Paletter,
	debugger Debugger,
	dayLimit float32,
) DayPrinter {
	return dayprint{
		output:   output,
		paletter: paletter,
		debug:    debugger,
		limit:    dayLimit,
	}
}

type dayprint struct {
	output   *termenv.Output
	paletter Paletter
	debug    Debugger
	limit    float32
}

func (dp dayprint) Print(d valueobject.Day) {
	weekday := GetWeekDayNumber(d.Date())

	format := "02  "
	if weekday == 7 { // print new line - next week of month
		format += "\n"
	}

	s := dp.output.String(d.Format(format))
	idx := int(math.Round(float64(d.Hours())))
	color := dp.paletter.Index(idx)
	hours := d.Hours()

	// apply style respecting to rules
	switch {
	case hours > 0:
		s = s.Background(dp.output.Color(color))
		if weekday >= 6 {
			dp.debug.Write(
				fmt.Sprintf("extra: %s +%.2fh/+%.2f", d.Format("Jan 02"), hours, d.Gross()))
		}
	case hours == 0 && weekday < 6 && d.IsPast():
		s = s.Background(dp.output.Color(color))
	case hours == 0 && weekday >= 6:
		s = s.Foreground(dp.output.Color(color))
	}

	if hours > dp.limit {
		dp.debug.Write(
			fmt.Sprintf("extra: %s +%.2fh", d.Format("2006-01-02"), hours-dp.limit))
	}

	fmt.Print(s)
}
