package apiclient

import (
	"time"

	"github.com/1buran/redmine"
	"github.com/1buran/workdiary/internal/domain/valueobject"
)

func NewRedmineApiClient(
	projectName string,
	url string,
	token string,
	userID string,
	hourlyRate float32,
	logging bool,
) ApiClient {
	return redmineApiClient{
		client: redmine.CreateApiClient(
			url, token, logging, redmine.TimeEntriesFilter{UserId: userID},
		),
		hourlyRate:  hourlyRate,
		projectName: projectName,
	}
}

type redmineApiClient struct {
	client      *redmine.ApiClient
	hourlyRate  float32 // employee rate
	projectName string
}

func (r redmineApiClient) List(d1, d2 time.Time) (<-chan valueobject.Day, <-chan error) {
	ch := make(chan valueobject.Day)
	er := make(chan error)
	r.client.StartDate = d1
	r.client.EndDate = d2
	go func() {
		defer close(ch)
		defer close(er)
		dataChan, errChan := redmine.Scroll[redmine.TimeEntries](r.client)
		for {
			select {
			case data, ok := <-dataChan:
				if ok { // data channel is open, perform action on the gotten item
					for _, item := range data.Items {
						day := valueobject.NewDay(item.SpentOn.Time)
						day.Track(r.hourlyRate, item.Hours)
						ch <- day
					}
					continue
				}
				return // data channel is closed, all data is transmitted, return to the main loop
			case err, ok := <-errChan:
				if ok { // err channel is open, perform action on the gotten error
					er <- err
				}
			}
		}
	}()
	return ch, er
}

func (r redmineApiClient) Project() string { return r.projectName }

// todo
func (r redmineApiClient) Track(date time.Time, hours float32) error { return nil }
