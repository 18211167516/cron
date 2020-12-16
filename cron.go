package cron

import (
	"sort"
	"sync"
	"time"
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

// byTime is a wrapper for sorting the entry array by time
// (with zero time at the end).
type byTime []*Task

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	// Two zero times should return false.
	// Otherwise, zero is "greater" than any other time.
	// (To sort it at the end of the list.)
	if s[i].Next.IsZero() {
		return false
	}
	if s[j].Next.IsZero() {
		return true
	}
	return s[i].Next.Before(s[j].Next)
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
	nextID    TaskID
}

// Info log
func (c *Cron) Info(value ...interface{}) {
	c.logger.Info(value...)
}

// Error log
func (c *Cron) Error(err error, value ...interface{}) {
	c.logger.Error(err, value...)
}

type FuncJob func()

func (f FuncJob) Run() { f() }

func (c *Cron) AddFunc(spec string, f func()) (TaskID, error) {
	return c.AddJob(spec, FuncJob(f))
}

// AddJob is The core method of adding Task entities
func (c *Cron) AddJob(spec string, job Job) (TaskID, error) {
	schedule, err := c.parser.Parse(spec)
	if err != nil {
		return 0, err
	}
	return c.Schedule(schedule, job), nil
}

// Schedule add
func (c *Cron) Schedule(schedule Schedule, cmd Job) TaskID {
	c.runningMu.Lock()
	defer c.runningMu.Unlock()
	c.nextID++
	task := &Task{
		ID:       c.nextID,
		Schedule: schedule,
		Job:      cmd,
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
		tasks:    nil,
		logger:   DefaultLogger,
		parser:   standardParser,
		location: time.Local,
	}
	return c
}

// now return location time
func (c *Cron) now() time.Time {
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

func (c *Cron) Run() {
	c.runningMu.Lock()
	if c.isRunning {
		c.runningMu.Unlock()
		return
	}
	c.isRunning = true
	c.runningMu.Unlock()
	c.run()
}

func (c *Cron) startJob(j Job) {
	go func() {
		j.Run()
	}()
}

// run
// step1 遍历全部task
func (c *Cron) run() {
	c.Info("Start")

	now := c.now()
	for _, task := range c.tasks {
		task.Next = task.Schedule.Next(now)
		c.Info("schedule", "now", now, "task", task.ID, "next", task.Next)
	}

	for {
		// Determine the next entry to run.
		sort.Sort(byTime(c.tasks))

		var timer *time.Timer
		if len(c.tasks) == 0 || c.tasks[0].Next.IsZero() {
			// If there are no entries yet, just sleep - it still handles new entries
			// and stop requests.
			timer = time.NewTimer(100000 * time.Hour)
		} else {
			//获取下一次timer时间
			timer = time.NewTimer(c.tasks[0].Next.Sub(now))
		}

		for {
			select {
			//定时器监听触发
			case now = <-timer.C:
				now = now.In(c.location)
				c.logger.Info("wake", "now", now)

				// Run every entry whose next time was less than now
				for _, e := range c.tasks {
					//下次执行时间未到或者是零点
					if e.Next.After(now) || e.Next.IsZero() {
						break
					}
					c.startJob(e.Job)
					e.Prev = e.Next
					e.Next = e.Schedule.Next(now)
					c.logger.Info("run", "now", now, "entry", e.ID, "next", e.Next)
				}

			//服务启动之后新增cron
			case newEntry := <-c.add:
				timer.Stop()
				now = c.now()
				newEntry.Next = newEntry.Schedule.Next(now)
				c.tasks = append(c.tasks, newEntry)
				c.logger.Info("added", "now", now, "entry", newEntry.ID, "next", newEntry.Next)

			case <-c.stop:
				timer.Stop()
				c.logger.Info("stop")
				return

			case id := <-c.remove:
				timer.Stop()
				now = c.now()
				c.removeTask(id)
				c.logger.Info("removed", "entry", id)
			}

			break
		}
	}
}

// removeTask 移除task
func (c *Cron) removeTask(id TaskID) {
	var entries []*Task
	for _, e := range c.tasks {
		if e.ID != id {
			entries = append(entries, e)
		}
	}
	c.tasks = entries
}
