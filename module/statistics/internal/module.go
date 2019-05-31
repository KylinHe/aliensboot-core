package internal

import (
	"github.com/KylinHe/aliensboot-core/module/base"
	"github.com/KylinHe/aliensboot-core/module/statistics/conf"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*base.Skeleton
}

func (m *Module) GetConfig() interface{} {
	return &conf.Config
}

func (m *Module) GetName() string {
	return "statistics"
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
}

func (m *Module) OnDestroy() {
}
