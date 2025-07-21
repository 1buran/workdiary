package usecase

import (
	"fmt"
	"time"

	"github.com/1buran/workdiary/internal/infrastructure/apiclient"
)

// Track time.
func Track(
	client apiclient.ApiClient, date time.Time,
	issue, activity string, hours float32, comment string,
) {
	if err := client.Track(date, issue, activity, hours, comment); err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
