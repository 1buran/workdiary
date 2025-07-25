package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/muesli/termenv"

	"github.com/1buran/workdiary/internal/domain/repository"
)

func NewCalendar(
	monthbegin time.Time,
	monthend time.Time,
	dayprint DayPrinter,
	repo repository.WorkdiaryRepository,
	expectedAmountColor, infactAmountColor, summaryColor string,
	debugger Debugger,
) Calendar {
	return calendar{
		begin:    monthbegin,
		end:      monthend,
		repo:     repo,
		dayprint: dayprint,
		output:   termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor)),

		debugger:            debugger,
		expectedAmountColor: expectedAmountColor,
		infactAmountColor:   infactAmountColor,
		summaryColor:        summaryColor,
	}
}

var (
	hr = strings.Repeat("-", 28) // horizontal rule

	captionFormat  = "Jan 2006"
	captionShift   = strings.Repeat(" ", len(hr)/2-len(captionFormat)+2)
	dayAbbrs       = "Mon Tue Wed Thu Fri Sat Sun"
	dayPlaceholder = strings.Repeat(" ", 4)
)

type calendar struct {
	begin    time.Time
	end      time.Time
	repo     repository.WorkdiaryRepository
	dayprint DayPrinter
	debugger Debugger
	output   *termenv.Output

	expectedAmountColor, infactAmountColor, summaryColor string
}

// Print header of calendar in format:
//
//		        Jan 2024
//
//	 ---------------------------
//	 Mon Tue Wed Thu Fri Sat Sun
//	 ---------------------------
//
// we used a new style, when the Monday is the first day of a week.
func (c calendar) PrintHeader() {
	fmt.Println(captionShift, c.output.String(c.begin.Format(captionFormat)).Bold())
	fmt.Println(hr)
	fmt.Println(c.output.String(dayAbbrs).Bold())
	fmt.Println(hr)
	fmt.Print(strings.Repeat(dayPlaceholder, GetWeekDayNumber(c.begin)-1))
}

// Print footer line.
//
//	---------------------------
//
// Nothing special, just a "horizontal rule".
func (c calendar) PrintFooter() {
	if GetWeekDayNumber(c.end) != 7 {
		fmt.Println() // print new line - next week of month
	}
	fmt.Println(hr)
}

// Print stats of employee in format:
//
//		Total  / In fact / Win or Lose
//	    ----------------------------
//		💰3312 / 3366.00 / 👍54.00
//
// win or lose means will take you higher or lower of expecting amount respectively.
// This is most motivation part of this program =).
func (c calendar) PrintSummary(rate, limit, total float32) {
	var s termenv.Style

	totalDays, past := c.CountWorkDays()
	expectedAmount := limit * rate * float32(totalDays)
	lossAmount := float32(past)*limit*rate - total

	s = c.output.String(fmt.Sprintf("💰%d", int(expectedAmount)))
	s = s.Foreground(c.output.Color(c.expectedAmountColor)).Bold()
	fmt.Print(s, " / ")

	s = c.output.String(fmt.Sprintf("%.2f", total))
	s = s.Foreground(c.output.Color(c.infactAmountColor)).Bold()
	fmt.Print(s, " / ")

	switch {
	case lossAmount < 0:
		s = c.output.String(fmt.Sprintf("👍%.2f", lossAmount*-1))
	case lossAmount > 0:
		s = c.output.String(fmt.Sprintf("😕%.2f", lossAmount))
	default:
		s = c.output.String(fmt.Sprintf("%.2f", lossAmount))
	}
	s = s.Foreground(c.output.Color(c.summaryColor)).Bold()
	fmt.Print(s, "\n")

	c.debugger.Writeln(strings.Repeat("-", 28))
	c.debugger.Writeln(fmt.Sprintf("Month expected hours: %.2f", float32(totalDays)*limit))
}

// Count work days of month.
func (c calendar) CountWorkDays() (total, past int) {
	today := time.Now()
	for d := c.begin; d.Before(c.end); d = d.AddDate(0, 0, 1) {
		if GetWeekDayNumber(d) < 6 {
			total++
			if d.Before(today) {
				past++
			}
		}
	}
	return
}

// Print employee work diary in calendar format.
func (c calendar) PrintDays() {
	for _, day := range c.repo.List() {
		c.dayprint.Print(day)
	}
}
