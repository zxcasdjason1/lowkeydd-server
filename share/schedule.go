package share

import (
	"time"
)

type Schedule struct {
	Event  *Emitter
	ticker *time.Ticker
	done   chan bool
	isStop bool
}

const (
	SCHE_UPDATE = "schedule@@update"
	SCHE_END    = "schedule@@end"
)

func NewSchedule(interval time.Duration) *Schedule {
	this := &Schedule{
		Event:  NewEmitter(),
		ticker: time.NewTicker(interval),
		done:   make(chan bool),
		isStop: false,
	}
	go func() {
		for {
			select {
			case <-this.ticker.C:
				this.Event.Emit(SCHE_UPDATE)
			case <-this.done:
				this.Event.Emit(SCHE_END)
				close(this.done)
				this.ticker.Stop()
				return
			}
		}
	}()
	return this
}

func (s *Schedule) Stop() {
	if !s.isStop {
		s.isStop = true
		s.done <- true
	}
}
