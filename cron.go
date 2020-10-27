package cron

import (
	"sync"
	"time"
	"fmt"
)

// TaskID is Task ID
type TaskID int

// Task is entity
type Task struct {
	ID       TaskID
	Schedule Schedule
	Job      Job
	Next     time.Time
	Prev     time.Time
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
	nextID 	  TaskID
}

// Info log
func (c *Cron) Info(value ...interface{}){
	c.logger.Info(value...)
}

// Error log
func (c *Cron) Error(err error,value ...interface{}){
	c.logger.Error(err,value...)
}

type FuncJob func()

func (f FuncJob) Run() { f() }

func (c *Cron) addFunc(spec string,f func())(TaskID,error){
	return c.AddJob(spec, FuncJob(f))
}

// AddJob is The core method of adding Task entities
func (c *Cron) AddJob(spec string,job Job) (TaskID, error){
	schedule, err := c.parser.Parse(spec)
	if err != nil {
		return 0, err
	}
	return c.Schedule(schedule,job),nil
}

// Schedule add 
func (c *Cron) Schedule(schedule Schedule, cmd Job) TaskID {
	c.runningMu.Lock()
	defer c.runningMu.Unlock()
	c.nextID++
	task := &Task{
		ID:         c.nextID,
		Schedule:   schedule,
		Job:        cmd,
	}
	if !c.isRunning {
		c.tasks = append(c.tasks, task)
	} else {
		c.add <- task
	}
	return task.ID
}

// New returns a new Cron job runner, modified by the given options
func New() *Cron {
	c := &Cron{
		tasks: nil,
		logger:DefaultLogger,
		parser:standardParser,
		location:time.Local,
	}
	return c
}

// now return location time
func (c *Cron) now() time.Time{
	return time.Now().In(c.location)
}

// Start the cron scheduler in its own goroutine, or no-op if already started.
func (c *Cron) Start() {
	c.runningMu.Lock()
	defer c.runningMu.Unlock()
	if c.isRunning {
		return
	}
	c.isRunning = true
	go c.run()
}

// run 
// step1 遍历全部task
func (c *Cron) run(){
	c.Info("Start")

	//now := c.now()
	//fmt.Println(now)
	fmt.Printf("%v,%d\n",c.tasks,len(c.tasks))
	for _, task := range c.tasks {
		c.Info("schedule", "now",task.ID)
	}

}
