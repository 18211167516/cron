package cron

import (
	"sync"
	"time"
)

// TaskID is Task ID
type TaskID int

// Task is entity
type Task struct {
	ID       TaskID
	Schedule Schedule
	Next     time.Time
	Prev     time.Time
	Job      Job
}

// Cron struct
type Cron struct {
	//
	tasks     []*Task
	stop      chan struct{}
	add       chan *Task
	remove    chan TaskID
	runningMu sync.Mutex
	isRunning bool
	logger    Logger
	location  *time.Location
	parser    ScheduleParser

}

func (c *Cron) Info(value ...interface{}){
	c.logger.Info(value...)
}

func (c *Cron) Error(err error,value ...interface{}){
	c.logger.Error(err,value...)
}

// AddJob is The core method of adding Task entities
func AddJob() {

}

// New returns a new Cron job runner, modified by the given options
func New() *Cron {

	c := &Cron{
		tasks: nil,
		logger:DefaultLogger,
	}

	return c
}

// Start the cron scheduler in its own goroutine, or no-op if already started.
func Start() {

}
