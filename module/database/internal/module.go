package internal

import (
	"github.com/KylinHe/aliensboot-core/module/base"
)

var (
	skeleton = base.NewSkeleton1(100000)
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*base.Skeleton
}

func (m *Module) GetName() string {
	return "database"
}

func (m *Module) GetConfig() interface{} {
	return nil
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton

}

func (m *Module) OnDestroy() {
}
