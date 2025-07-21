package apiclient

import (
	"math/rand/v2"
	"time"

	"github.com/1buran/workdiary/internal/domain/valueobject"
)

type demoApiClient struct{}

func (d demoApiClient) List(d1, d2 time.Time) (<-chan valueobject.Day, <-chan error) {
	ch := make(chan valueobject.Day)
	er := make(chan error)

	go func() {
		defer close(ch)
		defer close(er)
		for d := d1; d.Before(d2); d = d.AddDate(0, 0, 1) {
			day := valueobject.NewDay(d)
			day.Track(10, rand.Float32()*10)
			ch <- day
		}
	}()
	return ch, er
}

func (d demoApiClient) Project() string { return "demo" }
func (d demoApiClient) Track(
	date time.Time, issue, activity string, hours float32, comment string,
) error {
	return nil
}

func NewDemoApiClient() ApiClient { return demoApiClient{} }
