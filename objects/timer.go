package objects

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Timer struct {
	TimeLeft  float32
	StartTime float32
}

func NewTimer(time float32, start bool, gameState *int) *Timer {
	t := &Timer{TimeLeft: time, StartTime: time}
	if start {
		go t.SelfUpdate(gameState)
	}
	return t
}

func (t *Timer) Update(gameState *int) {
	if *gameState == 0 { // if game is playing
		t.TimeLeft -= rl.GetFrameTime()
	}
}

func (t *Timer) Reset() {
	t.TimeLeft = t.StartTime
}

func (t *Timer) IsDone() bool {
	return t.TimeLeft <= 0
}

func (t *Timer) SelfUpdate(gameState *int) {
	for {
		t.Update(gameState)
		time.Sleep(1 * time.Millisecond)
	}
}
