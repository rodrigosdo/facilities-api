package domain

import (
	"time"
)

type Shift struct {
	End      time.Time
	Facility Facility
	ID       int64
	Start    time.Time
}

type Shifts []Shift
