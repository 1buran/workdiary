package valueobject

import (
	"time"
)

// A Day holds the information about the work day: date and spent hours.
type Day struct {
	ts         time.Time
	spenthours float32
	gross      float32
	hourlyrate float32
}

func (d *Day) Track(rate, hours float32) {
	d.spenthours = hours
	d.hourlyrate = rate
	d.gross = hours * rate
}

func (d Day) Hours() float32         { return d.spenthours }
func (d Day) Gross() float32         { return d.gross }
func (d Day) Rate() float32          { return d.hourlyrate }
func (d Day) Date() time.Time        { return d.ts }
func (d Day) Format(f string) string { return d.ts.Format(f) }
func (d Day) IsPast() bool           { return d.ts.Before(time.Now()) }

func NewDay(d time.Time) Day { return Day{ts: d} }
func NewDayTracked(d time.Time, r, h, g float32) Day {
	return Day{ts: d, hourlyrate: r, spenthours: h, gross: g}
}
