package cron

// Job is an interface for submitted cron jobs.
type Job interface {
	Run()
}
