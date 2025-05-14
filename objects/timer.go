package objects

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Timer struct {
	TimeLeft  float32
	StartTime float32
}

func NewTimer(time float32, start bool) *Timer {
	t := &Timer{TimeLeft: time, StartTime: time}
	if start {
		go t.SelfUpdate()
	}
	return t
}

func (t *Timer) Update() {
	t.TimeLeft -= rl.GetFrameTime()
}

func (t *Timer) Reset() {
	t.TimeLeft = t.StartTime
}

func (t *Timer) IsDone() bool {
	return t.TimeLeft <= 0
}

func (t *Timer) SelfUpdate() {
	for {
		t.Update()
		time.Sleep(1 * time.Millisecond)
	}
}
