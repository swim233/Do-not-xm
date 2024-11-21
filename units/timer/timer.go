package timer

import (
	"fmt"
	"time"
)

var Time = new(int)

func CheckSleep() bool {
	if *Time <= 0 {

		return true
	} else {
		return false
	}
}

func Timer() {
	*Time = 0
	for {
		if *Time >= 0 {
			*Time = *Time - 1
			time.Sleep(1 * time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func Calculation(Time *int) string {
	var Hours = new(int)
	var Mins = new(int)
	*Mins = *Time / 60
	sec := *Time % 60
	*Hours = 0
	if *Mins >= 60 {
		*Hours = *Mins / 60
		*Mins = *Mins % 60
	}
	str := fmt.Sprintf("%d 小时 %d 分钟 %d 秒", *Hours, *Mins, sec)
	return str
}
