package core

import (
	"github.com/KylinHe/aliensboot-core/common/util"
	"time"
)

const (
	minTimerInterval = time.Millisecond * 10
)

type EntityTimerID int

func (this EntityTimerID) IsNil() bool {
	return this == 0
}

type entityTimerInfo struct {
	FireTime       time.Time
	RepeatInterval time.Duration
	Method         string
	Args           []interface{}
	Repeat         bool
	rawTimer       *util.Timer
}

// CancelTimer cancels the Callback / Timer
func (e *Entity) CancelTimer(tid EntityTimerID) {
	if tid.IsNil() {
		return
	}
	timerInfo := e.timers[tid]
	if timerInfo == nil {
		return
	}
	delete(e.timers, tid)
	e.cancelRawTimer(timerInfo.rawTimer)
}

func (e *Entity) cancelRawTimer(t *util.Timer) {
	delete(e.rawTimers, t)
	t.Cancel()
}

func (e *Entity) AddTimer(d time.Duration, repeat bool, method string, args ...interface{}) EntityTimerID {
	if d < minTimerInterval { // minimal interval for repeat timer
		d = minTimerInterval
	}
	tid := e.genTimerId()
	now := time.Now()
	info := &entityTimerInfo{
		FireTime:       now.Add(d),
		RepeatInterval: d,
		Method:         method,
		Args:           args,
		Repeat:         repeat,
	}
	e.timers[tid] = info
	info.rawTimer = e.addRawTimer(d, func() {
		e.triggerTimer(tid, true)
	})
	return tid
}

func (e *Entity) clearRawTimers() {
	for t := range e.rawTimers {
		t.Cancel()
	}
	e.rawTimers = nil
}

func (e *Entity) triggerTimer(tid EntityTimerID, isRepeat bool) {
	timerInfo := e.timers[tid] // should never be nil
	if timerInfo == nil {
		return
	}
	if !timerInfo.Repeat {
		delete(e.timers, tid)
	} else {
		if !isRepeat {
			timerInfo.rawTimer = e.addRawTimer(timerInfo.RepeatInterval, func() {
				e.triggerTimer(tid, true)
			})
		}
		now := time.Now()
		timerInfo.FireTime = now.Add(timerInfo.RepeatInterval)
	}

	e.OnCallFromLocal(timerInfo.Method, timerInfo.Args)
}

func (e *Entity) addRawTimer(d time.Duration, cb util.CallbackFunc) *util.Timer {
	t := EntityManager.addTimer(d, cb)
	e.rawTimers[t] = struct{}{}
	return t
}

func (e *Entity) addRawCallback(d time.Duration, cb util.CallbackFunc) *util.Timer {
	var t *util.Timer
	t = EntityManager.addCallback(d, func() {
		delete(e.rawTimers, t)
		cb()
	})
	e.rawTimers[t] = struct{}{}
	return t
}

func (e *Entity) genTimerId() EntityTimerID {
	e.lastTimerId += 1
	tid := e.lastTimerId
	return tid
}
