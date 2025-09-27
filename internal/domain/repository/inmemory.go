package repository

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/1buran/workdiary/internal/domain/valueobject"
)

// In memory repo holds the value objects in memory.
type inmemory struct {
	days []valueobject.Day

	sync.RWMutex
}

func (i *inmemory) List() (days []valueobject.Day) {
	i.RLock()
	defer i.RUnlock()

	return i.days
}

func (i *inmemory) Add(d valueobject.Day) {
	i.Lock()
	defer i.Unlock()

	i.days = append(i.days, d)
}

func (i *inmemory) Compact() {
	i.Lock()
	defer i.Unlock()

	comments := make(map[string][]string)
	diary := make(map[string][]float32)
	for _, d := range i.days {
		if v, ok := diary[d.Format(time.DateOnly)]; ok {
			v[1] += d.Hours()
			v[2] += d.Gross()
			v[0] = v[2] / v[1]
			comments[d.Format(time.DateOnly)] = append(
				comments[d.Format(time.DateOnly)], d.Comments())
		} else {
			diary[d.Format(time.DateOnly)] = []float32{d.Rate(), d.Hours(), d.Gross()}
			comments[d.Format(time.DateOnly)] = []string{d.Comments()}
		}
	}

	var dates []string
	for k := range diary {
		dates = append(dates, k)
	}
	sort.Strings(dates)

	var ndays []valueobject.Day
	for _, k := range dates {
		t, _ := time.Parse(time.DateOnly, k)
		v := valueobject.NewDayTracked(t, diary[k][0], diary[k][1], diary[k][2], strings.Join(comments[k], ";"))
		ndays = append(ndays, v)
	}

	i.days = ndays
}

func (i *inmemory) TotalHours() (total float32) {
	for _, day := range i.List() {
		total += day.Hours()
	}
	return
}

func (i *inmemory) TotalAmount() (total float32) {
	for _, day := range i.List() {
		total += day.Gross()
	}
	return
}

func (i *inmemory) MaxDayHours() (max float32) {
	for _, day := range i.List() {
		dh := day.Hours()
		if dh > max {
			max = dh
		}
	}
	return
}

func NewInMemoryRepository() WorkdiaryRepository {
	return &inmemory{}
}
