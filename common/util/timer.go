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
	"fmt"
	"os"
	"runtime/debug"
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
type CallbackFunc func()

// Add a callback which will be called after specified duration
func (manager *TimerManager) AddCallback(d time.Duration, callback CallbackFunc) *Timer {
	t := &Timer{
		fireTime: time.Now().Add(d),
		interval: d,
		callback: callback,
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

	t := &Timer{
		fireTime: time.Now().Add(d),
		interval: d,
		callback: callback,
		repeat:   true,
	}
	t.addSeq = manager.nextAddSeq
	manager.nextAddSeq += 1

	heap.Push(manager._TimerHeap, t)
	return t
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
		runCallback(callback)
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

func runCallback(callback CallbackFunc) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Callback %v paniced: %v\n", callback, err)
			debug.PrintStack()
		}
	}()
	callback()
}
