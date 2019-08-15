package core

import (
	"github.com/KylinHe/aliensboot-core/database"
	"reflect"
)

type DataOp uint8

const (
	OpInsert DataOp = iota
	OpUpdate
	OpDelete
)


type DataCacheConfig struct {
	UpdateInterval int64 //缓存队列写入数据库的间隔时间 单位 秒
	MaxCacheSize   int   //缓存队列的最大值，超出就需要开始写数据库了
	MaxWriteSize   int   //缓存单次写的最大值，超出的留到下一次时间写数据库
}

type DataOperation struct {
	op DataOp
	data interface{}
}

type DataCacheManager struct {
	caches map[database.IDatabaseHandler]*DatabaseCache
}

// 设置缓存策略
func (manager *DataCacheManager) RegisterStrategy(handler database.IDatabaseHandler, config DataCacheConfig) {
	manager.caches[handler] = &DatabaseCache{DataCacheConfig:config, caches:make(map[reflect.Type]*CollectionCache)}
}

func (manager *DataCacheManager) OpData(handler database.IDatabaseHandler) {

}

//所有集合的缓存
type DatabaseCache struct {
	DataCacheConfig
	caches map[reflect.Type]*CollectionCache
}

//所有集合数据的缓存  id - value
type CollectionCache struct {
	caches map[interface{}]DataOperation
}
