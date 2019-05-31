package task

import "time"

//--------------------
type ITask interface {
	Start()
	Close()
	IsRunning() bool
}

//-----------------------
type Runnable interface {
	execute()
}

//------------------------
type TimerTask struct {
	Ticker  *time.Ticker
	running bool
	Runnable
}

func (this *TimerTask) Start(executors ...func()) {
	go func() {
		this.running = true
		for _ = range this.Ticker.C {
			for _, executor := range executors {
				executor()
			}
		}
		this.running = false
	}()
}

func (this *TimerTask) Close() {
	this.Ticker.Stop()
}

func (this *TimerTask) isRunning() bool {
	return this.running
}
