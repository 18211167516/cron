package main

import (
	"fmt"
	"github.com/18211167516/cron"
)

func main(){
	c := cron.New()
	c.AddFunc("*",func(){
		fmt.Println("111111")
	})

	c.AddFunc("*",func(){
		fmt.Println("222222")
	})

	c.Start()
	select {} 
}