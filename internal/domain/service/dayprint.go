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
	showComments bool,
) DayPrinter {
	return dayprint{
		output:       output,
		paletter:     paletter,
		debug:        debugger,
		limit:        dayLimit,
		showComments: showComments,
	}
}

type dayprint struct {
	output       *termenv.Output
	paletter     Paletter
	debug        Debugger
	limit        float32
	showComments bool
}

func (dp dayprint) Print(d valueobject.Day) {
	weekday := GetWeekDayNumber(d.Date())

	format := "02  "
	if weekday == 7 { // print new line - next week of month
		format += "\n"
	}

	s := dp.output.String(d.Format(format))
	idx := int(math.Round(float64(d.Hours())))
	color, fgColor := dp.paletter.Index(idx)
	hours := d.Hours()

	s = s.Foreground(dp.output.Color(fgColor))

	if weekday >= 6 {
		dp.debug.Write(dp.output.String(d.Format("Mon Jan 02 ")).Faint())
	} else {
		dp.debug.Write(d.Format("Mon Jan 02 "))
	}

	// apply style respecting to rules
	switch {
	case hours > 0:
		s = s.Background(dp.output.Color(color))
		if weekday >= 6 {
			dp.debug.Write(
				fmt.Sprintf("extra: +%.2fh/+%.2f ", hours, d.Gross()))
		} else {
			dp.debug.Write(
				fmt.Sprintf("woday: %.2fh/%.2f ", hours, d.Gross()))
		}
	case hours == 0 && weekday < 6 && d.IsPast():
		s = s.Background(dp.output.Color(color))
		dp.debug.Write("--")
	case hours == 0 && weekday >= 6:
		s = s.Foreground(dp.output.Color(color))
	}

	if hours > dp.limit {
		extra := hours - dp.limit
		dp.debug.Write(fmt.Sprintf("extra: +%.2fh/+%.2f ", extra, extra*d.Rate()))
	}

	if hours > 0 && dp.showComments {
		dp.debug.Write(d.Comments())
	}

	dp.debug.Writeln()

	fmt.Print(s)
}
