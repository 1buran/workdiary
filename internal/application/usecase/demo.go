package usecase

import (
	"fmt"
	"os"
	"time"

	"github.com/muesli/termenv"

	"github.com/1buran/workdiary/internal/infrastructure/apiclient"
)

func Demo(themes [][3]string) { // format of array: [name, startColor, endColor]
	clients := []apiclient.ApiClient{apiclient.NewDemoApiClient()}
	today := time.Now()
	year := today.Year()

	output := termenv.NewOutput(os.Stdout)
	output.HideCursor()

	total := len(themes)
	for i := range 12 {
		idx := i % total
		theme := themes[idx]
		monthbegin := time.Date(year, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
		monthend := monthbegin.AddDate(0, 1, 0).Add(-time.Nanosecond)

		Show(
			output, clients, monthbegin, monthend,
			theme[1], theme[2], "#ff9ff3", "#4cd137", "#fd79a8",
			false,
		)
		fmt.Println(
			output.String("\ntheme>").Faint().String(),
			output.String(theme[0]).Foreground(termenv.ANSIBrightMagenta),
		)
		time.Sleep(400 * time.Millisecond)
		output.ClearScreen()
	}
	output.ShowCursor()
}
