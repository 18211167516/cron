package cron

import (
	"fmt"
)

func Example_cron(){
	c := New()
	c.addFunc("*",func(){
		fmt.Println("111111")
	})

	c.addFunc("*",func(){
		fmt.Println("222222")
	})

	c.Start()
	//Output:111
}