package apiclient

import (
	"time"

	"github.com/1buran/workdiary/internal/domain/valueobject"
)

type ApiClient interface {
	List(d1, d2 time.Time) (<-chan valueobject.Day, <-chan error)
	Track(date time.Time, issue, activity string, hours float32, comment string) error
	Project() string
}
