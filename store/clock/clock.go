package clock

import "time"

type Clocker interface {
	Now() time.Time
}

type RealClocker struct{}

func (r RealClocker) Now() time.Time {
	return time.Now()
}

type FixedClocker struct{}

func (rc FixedClocker) Now() time.Time {
	return time.Date(2023, 2, 22, 12, 34, 56, 0, time.UTC)
}
