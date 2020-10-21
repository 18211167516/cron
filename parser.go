package cron

import "time"

// Schedule describes a job's duty cycle.
type Schedule interface {
	// Next returns the next activation time, later than the given time.
	// Next is invoked initially, and then each time the job is run.
	Next(time.Time) time.Time
}

// ScheduleParser is Schedule Parser interface
type ScheduleParser interface {
	Parse(spec string) (Schedule, error)
}
