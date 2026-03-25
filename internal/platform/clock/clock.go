package clock

import "time"

type Clock interface {
	Now() time.Time
}

type SystemClock struct{}

func New() SystemClock {
	return SystemClock{}
}

func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}
