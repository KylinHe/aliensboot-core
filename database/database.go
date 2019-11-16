package database

import (
	"github.com/KylinHe/aliensboot-core/config"
	"github.com/KylinHe/aliensboot-core/database/dbconfig"
)

//数据库抽象层 适配其他数据库
type IDatabase interface {
	Init(config config.DBConfig) error //初始化数据库
	//auth(username string, password string)           //登录信息
	Close()                       //关闭数据库
	GetHandler() IDatabaseHandler //获取数据库处理类
}

type IDatabaseFactory interface {
	create() IDatabase
}

type IData interface {
	GetDataId() interface{} //获取数据id 不能为指针
}

type Authority struct {
	Username string
	Password string
}

//数据库handler
type IDatabaseHandler interface {
	GetTableMeta(data IData) (*dbconfig.TableMeta, error)
	//GetTableName(data IData) (string, error)
	EnsureTable(name string, data IData) error //确保表存在
	//GetID(data IData) interface{}
	EnsureIndex(name string, key []string, unique bool) error                                                       //确保索引

	//Related(data IData, relateData interface{}, relateTableName string, relateKey string) error //创建依赖关系
	//GenId(data IData) (int32, error)
	//InsertWithoutID(data IData) error
	//GenTimestampId(data IData) (int64, error)
	Insert(data IData) error
	InsertMulti(data []interface{}) error  //插入多条数据
	QueryAll(data IData, result interface{}) error
	QueryAllLimit(data IData, result interface{}, limit int, callback func(interface{}) bool) error
	QueryAllConditionLimit(data IData, condition string, value interface{}, result interface{}, limit int, callback func(interface{}) bool) error
	QueryAllConditionsLimit(data IData, conditions map[string]interface{}, result interface{}, limit int, sort ...string) error
	QueryAllConditionSkipLimit(data IData, condition string, value interface{}, result interface{}, skip int, limit int, sort ...string) error
	QueryAllConditionsSkipLimit(data IData, conditions map[string]interface{}, result interface{}, skip int, limit int, sort ...string) error
	QueryAllCondition(data IData, condition string, value interface{}, result interface{}) error
	QueryAllConditions(data IData, conditions map[string]interface{}, result interface{}) error
	QueryConditionCount(data IData, condition string, value interface{}) (int, error)
	QueryConditionsCount(data IData, query interface{}) (int, error)
	PipeAllConditions(data IData, pipeline interface{}, result interface{}) error
	QueryOne(data IData) error
	QueryOneCondition(data IData, condition string, value interface{}) error
	QueryOneConditions (data IData, conditions map[string]interface{}) error
	IDExist(data IData) (bool, error)
	DeleteOne(data IData) error
	DeleteOneCondition(data IData, selector interface{}) error
	DeleteAllCondition(data IData, selector interface{}) error
	UpdateOne(data IData) error
	ForceUpdateOne(data IData) error //强制更新。不存在就插入
	UpdateOneCondition(data IData, selector interface{}, update interface{}) error
}

//type interface{} interface {
//	//Name() string       //数据名称
//	//GetID() interface{} //获取数据索引ID
//}

type IRelatedData interface {
	RelateLoad(handler IDatabaseHandler) //注入关联数据
	RelateSave(handler IDatabaseHandler) //保存关联数据
}
