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

type Authority struct {
	Username string
	Password string
}

//数据库handler
type IDatabaseHandler interface {
	GetTableMeta(data interface{}) (*dbconfig.TableMeta, error)
	//GetTableName(data interface{}) (string, error)
	EnsureTable(name string, data interface{}) error //确保表存在
	//GetID(data interface{}) interface{}
	EnsureIndex(name string, key []string, unique bool) error                                                       //确保索引
	Related(data interface{}, relateData interface{}, relateTableName string, relateKey string) error //创建依赖关系
	//GenId(data interface{}) (int32, error)
	//InsertWithoutID(data interface{}) error
	//GenTimestampId(data interface{}) (int64, error)
	Insert(data interface{}) error
	QueryAll(data interface{}, result interface{}) error
	QueryAllLimit(data interface{}, result interface{}, limit int, callback func(interface{}) bool) error
	QueryAllConditionLimit(data interface{}, condition string, value interface{}, result interface{}, limit int, callback func(interface{}) bool) error
	QueryAllConditionsLimit(data interface{}, conditions map[string]interface{}, result interface{}, limit int, sort ...string) error
	QueryAllConditionSkipLimit(data interface{}, condition string, value interface{}, result interface{}, skip int, limit int, sort ...string) error
	QueryAllConditionsSkipLimit(data interface{}, conditions map[string]interface{}, result interface{}, skip int, limit int, sort ...string) error
	QueryAllCondition(data interface{}, condition string, value interface{}, result interface{}) error
	QueryAllConditions(data interface{}, conditions map[string]interface{}, result interface{}) error
	QueryConditionCount(data interface{}, condition string, value interface{}) (int, error)
	QueryConditionsCount(data interface{}, query interface{}) (int, error)
	PipeAllConditions(data interface{}, pipeline interface{}, result interface{}) error
	QueryOne(data interface{}) error
	QueryOneCondition(data interface{}, condition string, value interface{}) error
	QueryOneConditions (data interface{}, conditions map[string]interface{}) error
	IDExist(data interface{}) (bool, error)
	DeleteOne(data interface{}) error
	DeleteOneCondition(data interface{}, selector interface{}) error
	DeleteAllCondition(data interface{}, selector interface{}) error
	UpdateOne(data interface{}) error
	ForceUpdateOne(data interface{}) error //强制更新。不存在就插入
	Update(collection string, selector interface{}, update interface{}) error
}

//type interface{} interface {
//	//Name() string       //数据名称
//	//GetID() interface{} //获取数据索引ID
//}

type IRelatedData interface {
	RelateLoad(handler IDatabaseHandler) //注入关联数据
	RelateSave(handler IDatabaseHandler) //保存关联数据
}
