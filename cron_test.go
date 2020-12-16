/*
 * @Description: 项目
 * @Date: 2020-10-26 17:21:08
 * @Author: baichonghua
 * @Version: V1.0.0
 * @LastEditors: baichonghua
 * @LastEditTime: 2020-12-15 09:48:40
 */
package cron

import (
	"fmt"
)

func Example_cron() {
	c := New()
	c.AddFunc("*", func() {
		fmt.Println("111111")
	})

	c.AddFunc("*", func() {
		fmt.Println("222222")
	})

	c.Start()
	//Output:111
}
