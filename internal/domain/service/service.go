package service

import (
	"github.com/muesli/termenv"

	"github.com/1buran/workdiary/internal/domain/valueobject"
)

type Paletter interface {
	Sprint(output *termenv.Output) string
	Index(i int) string
}

type Debugger interface {
	Write(a ...any)
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
