package prometheus

import "time"

type TimeSeries struct {
	Labels []Label
	Sample Sample
	ID     int
}

type Label struct {
	Name  string
	Value string
}

type Sample struct {
	Time  time.Time
	Value float64
}
