package usecase

import (
	"fmt"
	"io"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/muesli/termenv"

	"github.com/1buran/workdiary/internal/domain/repository"
	"github.com/1buran/workdiary/internal/domain/service"
	"github.com/1buran/workdiary/internal/domain/valueobject"
	"github.com/1buran/workdiary/internal/infrastructure/apiclient"
)

// Show calendar, color of days depends on spent hours.
func Show(
	out io.Writer,
	clients []apiclient.ApiClient,
	d1, d2 time.Time,
	color1, color2 string,
	expectedAmountColor, infactAmountColor, summaryColor string,
	debug bool,
) {
	repo := repository.NewInMemoryRepository()

	// run all configured clients, gathering stats from Redmine, Gitlab etc
	var (
		agents   sync.WaitGroup
		projects []string
	)
	for _, client := range clients {
		projects = append(projects, client.Project())
		agents.Add(1)
		go func() {
			defer agents.Done()
			days, err := client.List(d1, d2)
			for {
				select {
				case day, ok := <-days:
					repo.Add(day)
					if !ok {
						return
					}
				case e, ok := <-err:
					if e != nil {
						fmt.Println(e)
					}
					if !ok {
						return
					}
				}
			}
		}()
	}
	agents.Wait()

	// fill the possible missed days: add all month days to repo.
	for d := d1; d.Before(d2); d = d.AddDate(0, 0, 1) {
		repo.Add(valueobject.NewDay(d))
	}

	repo.Compact() // compact the data: merge day stats from different sources

	// adaptive number of color gradients
	gradients := int(math.Ceil(float64(repo.MaxDayHours())))
	if gradients == 0 { // the diary is empty, full palette is no needed
		gradients = 2
	}

	output := termenv.NewOutput(out, termenv.WithProfile(termenv.TrueColor))
	paletter := service.NewPaletter(color1, color2, gradients)

	// show extra data for debug: rate, limit, expected hours in month
	rate := repo.TotalAmount() / repo.TotalHours()
	debugger := service.NewDebugger(debug)
	debugger.Write(strings.Repeat("-", 27))
	debugger.Write(paletter.Sprint(output))
	debugger.Write(fmt.Sprintf("projects: %q", projects))
	debugger.Write(fmt.Sprintf("rate: %.2f, hours: %.2f", rate, repo.TotalHours()))
	defer debugger.Read()

	limit := float32(8) // base daily limit of working hours

	dayprinter := service.NewDayPrinter(output, paletter, debugger, limit)
	cal := service.NewCalendar(d1, d2, dayprinter, repo, expectedAmountColor, infactAmountColor, summaryColor, debugger)

	cal.PrintHeader()
	cal.PrintDays()
	cal.PrintFooter()
	cal.PrintSummary(rate, limit, repo.TotalAmount())
}
