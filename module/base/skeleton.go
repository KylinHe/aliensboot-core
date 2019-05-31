/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/5/9
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package base

import (
	"github.com/KylinHe/aliensboot-core/chanrpc"
	"github.com/KylinHe/aliensboot-core/console"
	"github.com/KylinHe/aliensboot-core/pool"
	"github.com/KylinHe/aliensboot-core/task"
	"time"
)

const (
	// skeleton conf
	GoLen              = 10000
	TimerDispatcherLen = 10000
	AsynCallLen        = 10000
	ChanRPCLen         = 10000

	TickInterval = time.Millisecond * 5 // server tick interval => affect timer resolution
)

func NewSkeleton() *Skeleton {
	skeleton := &Skeleton{
		GoLen:              GoLen,
		TimerDispatcherLen: TimerDispatcherLen,
		AsynCallLen:        AsynCallLen,
		ChanRPCServer:      chanrpc.NewServer(ChanRPCLen),
		ticker:             time.Tick(TickInterval),
	}
	skeleton.Init()
	return skeleton
}

type Skeleton struct {
	GoLen              int
	TimerDispatcherLen int
	AsynCallLen        int
	ChanRPCServer      *chanrpc.Server
	g                  *pool.Go
	dispatcher         *task.Dispatcher
	client             *chanrpc.Client
	server             *chanrpc.Server
	commandServer      *chanrpc.Server

	ticker <-chan time.Time
	tick   func()
}

func (s *Skeleton) Init() {
	if s.GoLen <= 0 {
		s.GoLen = 0
	}
	if s.TimerDispatcherLen <= 0 {
		s.TimerDispatcherLen = 0
	}
	if s.AsynCallLen <= 0 {
		s.AsynCallLen = 0
	}

	s.g = pool.New(s.GoLen)
	s.dispatcher = task.NewDispatcher(s.TimerDispatcherLen)
	s.client = chanrpc.NewClient(s.AsynCallLen)
	s.server = s.ChanRPCServer

	if s.server == nil {
		s.server = chanrpc.NewServer(0)
	}
	s.commandServer = chanrpc.NewServer(0)
}

func (s *Skeleton) SetTick(tick func()) {
	s.tick = tick
}

func (s *Skeleton) Run(closeSig chan bool) {
	for {
		select {
		case <-closeSig:
			s.commandServer.Close()
			s.server.Close()

			for !s.g.Idle() || !s.client.Idle() {
				s.g.Close()
				s.client.Close()
			}
			return
		case ri := <-s.client.ChanAsynRet:
			s.client.Cb(ri)
		case ci := <-s.server.ChanCall:
			s.server.Exec(ci)
		case ci := <-s.commandServer.ChanCall:
			s.commandServer.Exec(ci)
		case cb := <-s.g.ChanCb:
			s.g.Cb(cb)
		case t := <-s.dispatcher.ChanTimer:
			t.Cb()
		case <-s.ticker:
			if s.tick != nil {
				s.tick()
			}
		}
	}
}

func (s *Skeleton) AfterFunc(d time.Duration, cb func()) *task.Timer {
	if s.TimerDispatcherLen == 0 {
		panic("invalid TimerDispatcherLen")
	}

	return s.dispatcher.AfterFunc(d, cb)
}

func (s *Skeleton) CronFunc(cronExpr *task.CronExpr, cb func()) *task.Cron {
	if s.TimerDispatcherLen == 0 {
		panic("invalid TimerDispatcherLen")
	}

	return s.dispatcher.CronFunc(cronExpr, cb)
}

func (s *Skeleton) Go(f func(), cb func()) {
	if s.GoLen == 0 {
		panic("invalid GoLen")
	}

	s.g.Go(f, cb)
}

func (s *Skeleton) NewLinearContext() *pool.LinearContext {
	if s.GoLen == 0 {
		panic("invalid GoLen")
	}

	return s.g.NewLinearContext()
}

func (s *Skeleton) AsynCall(server *chanrpc.Server, id interface{}, args ...interface{}) {
	if s.AsynCallLen == 0 {
		panic("invalid AsynCallLen")
	}

	s.client.Attach(server)
	s.client.AsynCall(id, args...)
}

func (s *Skeleton) RegisterChanRPC(id interface{}, f interface{}) {
	if s.ChanRPCServer == nil {
		panic("invalid ChanRPCServer")
	}

	s.server.Register(id, f)
}

func (s *Skeleton) RegisterCommand(name string, help string, f interface{}) {
	console.Register(name, help, f, s.commandServer)
}
