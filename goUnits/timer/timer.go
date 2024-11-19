package timer

import (
	"learn/goUnits/bot"
	"math/rand/v2"
	"time"
)

var Time = new(int)

func CheckSleep() bool {
	if *Time <= 0 {
		*Time = (rand.IntN(bot.BotConfig.RandomCD) + bot.BotConfig.StaticCD)
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
