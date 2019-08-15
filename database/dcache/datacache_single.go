package dcache

import (
	"github.com/KylinHe/aliensboot-core/database"
	"github.com/KylinHe/aliensboot-core/dispatch/lpc"
)

func NewSingleDataCache() DataCache {
	return &SingleDataCache{}
}

type SingleDataCache struct {
	op DataOp
	data database.IData
}

func (cache SingleDataCache) OpData(op DataOp, data database.IData) {
	if op == OpDelete || cache.op != OpInsert {
		cache.op = op
	}
	cache.data = data
}

// 导出数据
func (cache SingleDataCache) Flush(dbHandler database.IDatabaseHandler) {
	if dbHandler == nil {
		return
	}
	if cache.data == nil {
		return
	}
	if cache.op == OpInsert {
		lpc.DBServiceProxy.Insert(cache.data, dbHandler)
	} else if cache.op == OpUpdate {
		lpc.DBServiceProxy.Update(cache.data, dbHandler)
	} else if cache.op == OpDelete {
		lpc.DBServiceProxy.Delete(cache.data, dbHandler)
	}
	cache.data = nil
}


