package usecase

import (
	"os"
	"time"

	"github.com/muesli/termenv"

	"github.com/1buran/workdiary/internal/infrastructure/apiclient"
)

func Demo() {
	clients := []apiclient.ApiClient{apiclient.NewDemoApiClient()}
	today := time.Now()
	year := today.Year()

	output := termenv.NewOutput(os.Stdout)
	output.HideCursor()
	for i := range 12 {
		monthbegin := time.Date(year, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
		monthend := monthbegin.AddDate(0, 1, 0).Add(-time.Nanosecond)

		Show(
			output, clients, monthbegin, monthend,
			"#a958ad", "#0d420d", "#ff9ff3", "#4cd137", "#fd79a8",
			false,
		)
		time.Sleep(400 * time.Millisecond)
		output.ClearScreen()
	}
	output.ShowCursor()
}
