package tracing

import (
	"github.com/KylinHe/aliensboot-core/module/tracing/internal"
)

var (
	Module  = new(internal.Module)
	ChanRPC = internal.ChanRPC
)
