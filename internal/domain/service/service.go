package service

import (
	"github.com/muesli/termenv"

	"github.com/1buran/workdiary/internal/domain/valueobject"
)

type Paletter interface {
	AddBackgroundColor(c string)
	AddForegroundColor(c string)
	Index(i int) (bgColor, fgColor string)
	Sprint(output *termenv.Output) string
}

type Debugger interface {
	Write(a ...any)
	Writeln(a ...any)
	Read()
}

type DayPrinter interface {
	Print(d valueobject.Day)
}

type Calendar interface {
	PrintHeader()
	PrintDays()
	PrintFooter()
	PrintSummary(rate, limit, total float32)
}
