package cron

import (
	"time"
)

type SpecSchedule struct {
	
}

func (s *SpecSchedule) Next(t time.Time) time.Time {
	return t.Add(1*time.Second)
}