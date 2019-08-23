package util

type TimerGroup map[interface{}]*Timer

func (t TimerGroup) Cancel() {
	for _, timer := range t {
		if timer != nil {
			timer.Cancel()
		}
	}
}

func (t TimerGroup) CancelById(id interface{}) {
	timer := t[id]
	if timer != nil {
		timer.Cancel()
	}
}

func (t TimerGroup) Add(id interface{}, timer *Timer) {
	t.CancelById(id)
	t[id] = timer
}