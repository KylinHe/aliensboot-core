package util

import "time"

type TimerType int32


func NewTimerContainer(manager *TimerManager) *TimerContainer {
	return &TimerContainer{
		manager: manager,
		groupMapping:  make(map[TimerType]TimerGroup),
		singleMapping: make(map[TimerType]*Timer),
		filter: make(map[TimerType]bool),
	}
}

type TimerContainer struct {
	manager *TimerManager

	groupMapping map[TimerType]TimerGroup

	singleMapping map[TimerType]*Timer

	filter map[TimerType]bool //定时器过滤
}

func(container *TimerContainer) AddFilter(timerTypes ...TimerType)  {
	for _, timerType := range timerTypes {
		container.filter[timerType] = true
	}
}

func(container *TimerContainer) ensureGetGroup(timerType TimerType) TimerGroup {
	group := container.groupMapping[timerType]
	if group == nil {
		group = make(TimerGroup)
		container.groupMapping[timerType] = group
	}
	return group
}

func(container *TimerContainer) CancelAll() {
	for _, group := range container.groupMapping {
		group.Cancel()
	}
	for _, timer := range container.singleMapping {
		if timer != nil {
			timer.Cancel()
		}
	}
}

func(container *TimerContainer) setGroupTimer(timerType TimerType, id interface{}, timer *Timer) {
	group := container.ensureGetGroup(timerType)
	group.Add(id, timer)
}

func(container *TimerContainer) CancelGroupTimer(timerType TimerType, id interface{}) {
	group := container.groupMapping[timerType]
	if group == nil {
		return
	}
	group.CancelById(id)
}

// 设置单例定时器、会取消上个定时器
func(container *TimerContainer) setSingleTimer(timerType TimerType, timer *Timer) {
	oldTimer := container.singleMapping[timerType]
	if oldTimer != nil {
		oldTimer.Cancel()
	}
	container.singleMapping[timerType] = timer
}

func(container *TimerContainer) CancelSingleTimer(timerType TimerType) {
	timer := container.singleMapping[timerType]
	if timer == nil {
		return
	}
	timer.Cancel()
	delete(container.singleMapping, timerType)
}

func(container *TimerContainer) SetSingleTimestampCallback(timerType TimerType, timestamp int64, callback CallbackFunc, param ...interface{}) {
	if container.filter[timerType] {
		return
	}
	timer := container.manager.AddTimestampCallback(timestamp, callback, param...)
	container.setSingleTimer(timerType, timer)
}

func(container *TimerContainer) SetSingleTimestampTimer(timerType TimerType, timestamp int64, duration time.Duration, callback CallbackFunc) {
	if container.filter[timerType] {
		return
	}
	timer := container.manager.AddTimestampTimer(timestamp, duration, callback)
	container.setSingleTimer(timerType, timer)
}

func(container *TimerContainer) SetGroupTimestampCallback(timerType TimerType, id interface{}, timestamp int64, callback CallbackFunc, param ...interface{}) {
	if container.filter[timerType] {
		return
	}
	timer := container.manager.AddTimestampCallback(timestamp, callback, param...)
	container.setGroupTimer(timerType, id, timer)
}

func(container *TimerContainer) SetGroupTimestampTimer(timerType TimerType, id interface{}, timestamp int64, duration time.Duration, callback CallbackFunc) {
	if container.filter[timerType] {
		return
	}
	timer := container.manager.AddTimestampTimer(timestamp, duration, callback)
	container.setGroupTimer(timerType, id, timer)
}


func(container *TimerContainer) SetSingleTimer(timerType TimerType, duration time.Duration, callback CallbackFunc) {
	if container.filter[timerType] {
		return
	}
	timer := container.manager.AddTimer(duration, callback)
	container.setSingleTimer(timerType, timer)
}


//func (container *TimerContainer) addTimer(duration time.Duration, callback CallbackFunc) *Timer {
//
//	return container.manager.AddTimer(duration, callback)
//}