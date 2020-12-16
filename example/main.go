/*
 * @Description: 项目
 * @Date: 2020-10-27 12:07:51
 * @Author: baichonghua
 * @Version: V1.0.0
 * @LastEditors: baichonghua
 * @LastEditTime: 2020-12-16 14:47:25
 */
package main

import (
	"fmt"

	"github.com/18211167516/cron"
)

func main() {
	c := cron.New()
	c.AddFunc("* * * * * *", func() {
		fmt.Println("111111")
	})

	c.Run()
	//select {}
}
