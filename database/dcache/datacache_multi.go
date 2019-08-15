package dcache

import (
	"github.com/KylinHe/aliensboot-core/database"
	"github.com/KylinHe/aliensboot-core/dispatch/lpc"
)

func NewMultiDataCache() DataCache {
	return &MultiDataCache{caches:make(map[interface{}]*DataOperation)}
}

type MultiDataCache struct {
	caches map[interface{}]*DataOperation
}

func (cache MultiDataCache) OpData(op DataOp, data database.IData) {
	id := data.GetDataId()
	oldDop := cache.caches[id]
	if oldDop == nil {
		cache.caches[id] = &DataOperation{op:op, data:data}
	} else {
		if op == OpDelete || oldDop.op != OpInsert {
			oldDop.op = op
		}
		oldDop.data = data
	}
}

// 导出数据
func (cache MultiDataCache) Flush(dbHandler database.IDatabaseHandler) {
	if dbHandler == nil {
		return
	}
	if len(cache.caches) == 0 {
		return
	}
	var insertData = make([]database.IData, 0)
	//异步入库
	for _, dop := range cache.caches {
		if dop.op == OpInsert {
			insertData = append(insertData, dop.data)
		} else if dop.op == OpUpdate {
			lpc.DBServiceProxy.Update(dop.data, dbHandler)
		} else if dop.op == OpDelete {
			lpc.DBServiceProxy.Delete(dop.data, dbHandler)
		}
	}
	// 目前只支持批量插入
	lpc.DBServiceProxy.InsertMulti(insertData, dbHandler)
	cache.caches = make(map[interface{}]*DataOperation)
}


