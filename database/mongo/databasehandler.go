package mongo

import (
	"github.com/KylinHe/aliensboot-core/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	//"strconv"
	"strings"
	//"time"
	"github.com/KylinHe/aliensboot-core/database/dbconfig"
	"github.com/pkg/errors"
)


const (
	IdStore        string = "_id"
	IncreaseIdBase int    = 100000
)


//获取表格名和id值
func (this *Database) GetTableMeta(data interface{}) (*dbconfig.TableMeta, error) {
	tableType := reflect.TypeOf(data)
	result, ok := this.tableMetas[tableType]
	if !ok {
		return result, errors.New("un expect db collection " + tableType.String())
	}
	return result, nil
}

//func (this *Database) GetID(data interface{}) interface{} {
//	tableMeta, err := this.GetTableMeta(data)
//	if err != nil {
//		return -1
//	}
//	return this.reflectID(data, tableMeta.IDName)
//}

func (this *Database) reflectID(data interface{}, idName string) interface{} {
	return reflect.ValueOf(data).Elem().FieldByName(idName).Interface()
}



//新增自增长键
//func (this *Database) EnsureCounter(data interface{}) {
//	this.validateConnection()
//	tableMeta := this.GetTableMeta(data)
//	this.dbContext.EnsureCounter(this.dbName, IdStore, tableMeta)
//}

//确保索引
//func (this *Database) EnsureUniqueIndex(data interface{}, key string) error {
//	this.validateConnection()
//	tableMeta, err := this.GetTableMeta(data)
//	if err != nil {
//		return err
//	}
//	return this.dbContext.EnsureUniqueIndex(this.dbName, tableMeta, []string{key})
//}

func (this *Database) EnsureTable(name string, data interface{}) error {
	//this.validateConnection()
	tableType := reflect.TypeOf(data)
	if tableType == nil || tableType.Kind() != reflect.Ptr {
		return errors.New("table data pointer required")
	}

	meta := &dbconfig.TableMeta{Name:name}
	dataType := tableType.Elem()
	for i:=0; i<dataType.NumField(); i++ {
		field := dataType.Field(i)
		uniqueValue := field.Tag.Get("unique")
		if strings.Contains(uniqueValue, "true") {
			key := field.Tag.Get("bson")
			if key != "" {
				err := this.dbContext.EnsureUniqueIndex(this.dbName, name, []string{key})
				if err != nil {
					return err
				}
			}
		} else if strings.Contains(uniqueValue, "false") {
			key := field.Tag.Get("bson")
			if key != "" {
				err := this.dbContext.EnsureIndex(this.dbName, name, []string{key})
				if err != nil {
					return err
				}
			}
		} else {
			idKey := field.Tag.Get("bson")
			if idKey == IdStore {
				meta.IDName = field.Name
				value := field.Tag.Get("gorm")
				if strings.Contains(value, "AUTO_INCREMENT") {
					meta.AutoIncrement = true
					err := this.dbContext.EnsureCounter(this.dbName, IdStore, name, IncreaseIdBase)
					if err != nil {
						log.Debugf("[%v] ensure count err : %v", this.dbName, err)
					}
				}
			}
		}
	}
	if meta.IDName == "" {
		return errors.New("bson:_id is not found in " + name + " tag",)
	}
	this.tableMetas[tableType] = meta
	return nil
}

func (this *Database) Related(data interface{}, relateData interface{}, relateTableName string, relateKey string) error {
	//mongo采用树形结构，不用建立关系
	return nil
}

func (this *Database) EnsureIndex(name string, key []string, unique bool) error {
	if unique {
		return this.dbContext.EnsureUniqueIndex(this.dbName, name, key)
	}
	return this.dbContext.EnsureIndex(this.dbName, name, key)
}

//自增长id
//func (this *Database) GenId(data interface{}) (int32, error) {
//	this.validateConnection()
//	tableMeta, err := this.GetTableMeta(data)
//	if err != nil {
//		return -1, err
//	}
//	newid, _ := this.dbContext.NextSeq(this.dbName, IdStore, tableMeta.Name)
//	newid += IncreaseIdBase
//	return int32(newid), nil
//}
//
//func (this *Database) InsertWithoutID(data interface{}) error {
//	this.validateConnection()
//	tableMeta, err := this.GetTableMeta(data)
//	if err != nil {
//		return err
//	}
//	return collection.Insert(data)
//}

//插入新数据
func (this *Database) Insert(data interface{}) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		if tableMeta.AutoIncrement {
			newId, err1 := this.dbContext.NextSeq(this.dbName, IdStore, tableMeta.Name)
			if err1 != nil {
				return err1
			}
			reflect.ValueOf(data).Elem().FieldByName(tableMeta.IDName).SetInt(newId)
		}

		return collection.Insert(data)
	})
}


func (this *Database) QueryAllLimit(data interface{}, result interface{}, limit int, callback func(interface{}) bool) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		skip := 0
		for {
			err := collection.Find(nil).Limit(limit).Skip(skip).All(result)
			if err != nil {
				return err
			}
			skip += limit
			if callback(result) {
				return nil
			}
		}
	})

}

//查询所有数据
func (this *Database) QueryAll(data interface{}, result interface{}) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		return collection.Find(nil).All(result)
	})
}

//查询单条记录
func (this *Database) QueryOne(data interface{}) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		return collection.FindId(this.reflectID(data, tableMeta.IDName)).One(data)
	})

}

func (this *Database) DeleteOne(data interface{}) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		return collection.RemoveId(this.reflectID(data, tableMeta.IDName))
	})
}

func (this *Database) DeleteOneCondition(data interface{}, selector interface{}) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		return collection.Remove(selector)
	})
}

func (this *Database) DeleteAllCondition(data interface{}, selector interface{}) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		_, err := collection.RemoveAll(selector)
		return err
	})
}

//查询单条记录
func (this *Database) IDExist(data interface{}) (bool, error) {
	return this.BoolRef(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) (bool, error) {
		count, err := collection.FindId(this.reflectID(data, tableMeta.IDName)).Count()
		return count != 0, err
	})

}


//按条件多条查询
func (this *Database) QueryAllConditionLimit(data interface{}, condition string, value interface{}, result interface{}, limit int, callback func(interface{}) bool) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		skip := 0
		for {
			err := collection.Find(bson.M{condition: value}).Limit(limit).Skip(skip).All(result)
			if err != nil {
				return err
			}
			skip += limit
			if callback(result) {
				return nil
			}
		}
	})

}

func (this *Database) QueryAllConditionsLimit(data interface{}, conditions map[string]interface{}, result interface{}, limit int, sort ...string) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		return collection.Find(conditions).Sort(sort...).Limit(limit).All(result)
	})
}

func (this *Database) QueryAllConditionSkipLimit(data interface{}, condition string, value interface{}, result interface{}, skip int, limit int, sort ...string) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		return collection.Find(bson.M{condition: value}).Sort(sort...).Limit(limit).Skip(skip).All(result)
	})
}

func (this *Database) QueryAllConditionsSkipLimit(data interface{}, conditions map[string]interface{}, result interface{}, skip int, limit int, sort ...string) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		return collection.Find(conditions).Sort(sort...).Limit(limit).Skip(skip).All(result)
	})
}

//按条件多条查询
func (this *Database) QueryAllCondition(data interface{}, condition string, value interface{}, result interface{}) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		return collection.Find(bson.M{condition: value}).All(result)
	})

}

func (this *Database) QueryAllConditions(data interface{}, conditions map[string]interface{}, result interface{}) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		return collection.Find(conditions).All(result)
	})

}

func (this *Database) QueryConditionCount(data interface{}, condition string, value interface{}) (int, error) {
	return this.IntRef(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) (int, error) {
		return collection.Find(bson.M{condition: value}).Count()
	})

}

func (this *Database) QueryConditionsCount(data interface{}, query interface{}) (int, error) {
	return this.IntRef(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) (int, error) {
		return collection.Find(query).Count()
	})

}

func (this *Database) PipeAllConditions(data interface{}, pipeline interface{}, result interface{}) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		return collection.Pipe(pipeline).All(result)
	})
}

func (this *Database) QueryOneConditions (data interface{}, conditions map[string]interface{}) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		return collection.Find(conditions).One(data)
	})
}

//按条件单条查询
func (this *Database) QueryOneCondition(data interface{}, condition string, value interface{}) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		return collection.Find(bson.M{condition: value}).One(data)
	})
}

//更新单条数据
func (this *Database) UpdateOne(data interface{}) error {
	return this.Ref(data, func(tableMeta *dbconfig.TableMeta, collection *mgo.Collection) error {
		return collection.UpdateId(this.reflectID(data, tableMeta.IDName), data)
	})
}

func (this *Database) ForceUpdateOne(data interface{}) error {
	result, err := this.IDExist(data)
	if err != nil {
		return err
	}
	if result {
		return this.UpdateOne(data)
	} else {
		return this.Insert(data)
	}
}


