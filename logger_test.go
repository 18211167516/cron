package cron

import (
	"time"
	"errors"
)
func Example_info(){
	c := New()
	c.Info("schedule","now",time.Now(),"entry","1","next","2")
	c.Error(errors.New("无此权限"),"schedule","now",time.Now(),"entry","1","next","2")
	// Output:123
}