/*******************************************************************************
* Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
* All rights reserved.
* Date:
*     2018/12/6
* Contributors:
*     aliens idea(xiamen) Corporation - initial API and implementation
*     jialin.he <kylinh@gmail.com>
*******************************************************************************/
package util

import (
	"container/heap"
	"github.com/KylinHe/aliensboot-core/exception"
	"time"
)

const (
	minInterval = 1 * time.Millisecond
)

var DefaultTimerManager *TimerManager = nil

func init() {
	DefaultTimerManager = NewTimerManager()
}

func NewTimerManager() *TimerManager {
	var timerHeap = &_TimerHeap{}
	heap.Init(timerHeap)
	return &TimerManager{
		_TimerHeap: timerHeap,
		nextAddSeq: 1,
	}
}

type TimerManager struct {
	*_TimerHeap
	nextAddSeq uint
}

type Timer struct {
	fireTime time.Time
	interval time.Duration
	callback CallbackFunc
	param []interface{}
	repeat   bool
	addSeq   uint
}

func (t *Timer) Cancel() {
	t.callback = nil
}

func (t *Timer) IsActive() bool {
	return t.callback != nil
}

type _TimerHeap struct {
	timers []*Timer
}

func (h *_TimerHeap) Len() int {
	return len(h.timers)
}

func (h *_TimerHeap) Less(i, j int) bool {
	//log.Println(h.timers[i].fireTime, h.timers[j].fireTime)
	t1, t2 := h.timers[i].fireTime, h.timers[j].fireTime
	if t1.Before(t2) {
		return true
	}

	if t1.After(t2) {
		return false
	}
	// t1 == t2, making sure Timer with same deadline is fired according to their add order
	return h.timers[i].addSeq < h.timers[j].addSeq
}

func (h *_TimerHeap) Swap(i, j int) {
	var tmp *Timer
	tmp = h.timers[i]
	h.timers[i] = h.timers[j]
	h.timers[j] = tmp
}

func (h *_TimerHeap) Push(x interface{}) {
	h.timers = append(h.timers, x.(*Timer))
}

func (h *_TimerHeap) Pop() (ret interface{}) {
	l := len(h.timers)
	h.timers, ret = h.timers[:l-1], h.timers[l-1]
	return
}

// Type of callback function
type CallbackFunc func(param []interface{})

func (manager *TimerManager) AddTimeCallback(fireTime time.Time, callback CallbackFunc, param ...interface{}) *Timer {
	t := &Timer{
		fireTime: fireTime,
		interval: fireTime.Sub(time.Now()),
		callback: callback,
		param: param,
		repeat:   false,
	}
	t.addSeq = manager.nextAddSeq // set addseq when locked
	manager.nextAddSeq += 1

	heap.Push(manager._TimerHeap, t)
	return t
}

func (manager *TimerManager) AddTimestampCallback(timestamp int64, callback CallbackFunc, param ...interface{}) *Timer {
	fireTime := time.Unix(timestamp, 0)
	return manager.AddTimeCallback(fireTime, callback, param...)
}

// Add a callback which will be called after specified duration
func (manager *TimerManager) AddCallback(d time.Duration, callback CallbackFunc, param ...interface{}) *Timer {
	t := &Timer{
		fireTime: time.Now().Add(d),
		interval: d,
		callback: callback,
		param:  param,
		repeat:   false,
	}
	t.addSeq = manager.nextAddSeq // set addseq when locked
	manager.nextAddSeq += 1

	heap.Push(manager._TimerHeap, t)
	return t
}

// Add a timer which calls callback periodly
func (manager *TimerManager) AddTimer(d time.Duration, callback CallbackFunc) *Timer {
	if d < minInterval {
		d = minInterval
	}
	fireTime := time.Now().Add(d)
	return manager.addtimer(fireTime,d,callback)
}

func (manager *TimerManager) AddTimer1(timestamp int64,d time.Duration, callback CallbackFunc) (*Timer,bool) {
	if timestamp < time.Now().Unix() {
		return nil,false
	}
	if d < minInterval {
		d = minInterval
	}
	fireTime := time.Unix(timestamp, 0)
	return manager.addtimer(fireTime,d,callback),true
}

// Tick once for timers
func (manager *TimerManager) Tick() {
	now := time.Now()
	//timerHeapLock.Lock()
	//defer timerHeapLock.Unlock()
	for {
		if manager.Len() <= 0 {
			break
		}

		nextFireTime := manager.timers[0].fireTime
		if nextFireTime.After(now) {
			break
		}

		t := heap.Pop(manager._TimerHeap).(*Timer)

		callback := t.callback
		if callback == nil {
			continue
		}

		if !t.repeat {
			t.callback = nil
		}
		// unlock the lock to run callback, because callback may add more callbacks / timers
		//timerHeapLock.Unlock()

		//callback(t.param)
		runCallback(callback, t.param)
		//timerHeapLock.Lock()

		if t.repeat {
			// add Timer back to heap
			t.fireTime = t.fireTime.Add(t.interval)
			if !t.fireTime.After(now) { // might happen when interval is very small
				t.fireTime = now.Add(t.interval)
			}
			t.addSeq = manager.nextAddSeq
			manager.nextAddSeq += 1
			heap.Push(manager._TimerHeap, t)
		}
	}

}

// Start the self-ticking routine, which ticks per tickInterval
func (manager *TimerManager) StartTicks(tickInterval time.Duration) {
	go manager.selfTickRoutine(tickInterval)
}

func (manager *TimerManager) selfTickRoutine(tickInterval time.Duration) {
	for {
		time.Sleep(tickInterval)
		manager.Tick()
	}
}

func runCallback(callback CallbackFunc, param []interface{}) {
	defer func() {
		if err := recover(); err != nil {
			exception.PrintStackDetail(err)
		}
	}()
	callback(param)
}

//
func (manager *TimerManager) addtimer(fireTime time.Time,d time.Duration, callback CallbackFunc) *Timer {
	t := &Timer{
		fireTime: fireTime,
		interval: d,
		callback: callback,
		repeat:   true,
	}
	t.addSeq = manager.nextAddSeq
	manager.nextAddSeq += 1

	heap.Push(manager._TimerHeap, t)
	return t
}
