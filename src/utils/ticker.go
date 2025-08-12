package utils

import (
	"math/rand/v2"
	"time"
)

type Timer struct {
	CurrentTime  time.Duration
	CoolDownTime time.Duration
}

func (t *Timer) ChangeTime(nStaticTime time.Duration, nRandomTime time.Duration) {
	t.CoolDownTime = nStaticTime + time.Duration(rand.Int64N(int64(nRandomTime.Seconds())))
	t.CurrentTime = 0
}

func (t *Timer) Timer() {
	t.CurrentTime = 0
	for {
		if t.CurrentTime > 0 {
			t.CurrentTime -= time.Second
		}
		time.Sleep(1 * time.Second)

	}
}

func (t *Timer) GetTime() time.Duration {
	return t.CurrentTime
}
