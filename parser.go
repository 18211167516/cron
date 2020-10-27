package cron

import "time"

// Schedule describes a job's duty cycle.
type Schedule interface {
	// Next returns the next activation time, later than the given time.
	// Next is invoked initially, and then each time the job is run.
	Next(time.Time) time.Time
}

type ParseOption int

// ScheduleParser is Schedule Parser interface
type ScheduleParser interface {
	Parse(spec string) (Schedule, error)
}

type Parser struct{
	options ParseOption
}

func NewParser(options ParseOption) Parser{
	return Parser{options}
}

func (P Parser) Parse(spec string) (Schedule, error){
	return &SpecSchedule{}, nil
}

var standardParser = NewParser(1)

