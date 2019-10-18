package util

type TimerGroup map[interface{}]*Timer

func (t TimerGroup) Cancel() {
	for _, timer := range t {
		if timer != nil {
			timer.Cancel()
		}
	}
	t = make(map[interface{}]*Timer)
}

func (t TimerGroup) CancelById(id interface{}) {
	timer := t[id]
	if timer != nil {
		timer.Cancel()
		delete(t, id)
	}
}

func (t TimerGroup) Get(id interface{}) *Timer{
	return t[id]
}

func (t TimerGroup) Add(id interface{}, timer *Timer) {
	t.CancelById(id)
	t[id] = timer
}