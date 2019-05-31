/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/3/29
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *
 * Description:
 *      operation for golang struct
 *******************************************************************************/
package redis

import (
	"errors"
	"fmt"
	"github.com/KylinHe/aliensboot-core/common/util"
	"reflect"
	"strings"
)

const (
	tag        = "rorm"
	ignoreFlag = "-"
)

func newHashKey(data interface{}, id interface{}) string {
	hashName := fmt.Sprintf("%v:%v", reflect.TypeOf(data), id)
	splstr := strings.Split(hashName, ".")
	if len(splstr) > 1 {
		hashName = splstr[1]
	}
	return hashName
}

func (this *RedisCacheClient) OExists(data interface{}, id interface{}) (bool, error) {
	key := newHashKey(data, id)
	return this.Exists(key)
}

func (this *RedisCacheClient) OGetBoolFieldByID(data interface{}, id interface{}, fieldName string) (bool, error) {
	key := newHashKey(data, id)
	return this.HGetBool(key, fieldName)
}

func (this *RedisCacheClient) OGetFieldByID(data interface{}, id interface{}, fieldName string) (string, error) {
	key := newHashKey(data, id)
	return this.HGet(key, fieldName)
}

func (this *RedisCacheClient) OSetFieldByID(data interface{}, id interface{}, fieldName string, value interface{}) error {
	key := newHashKey(data, id)
	return this.HSet(key, fieldName, value)
}

//获取多个字段
func (this *RedisCacheClient) OGetFieldsByID(data interface{}, id interface{}, fieldNames ...interface{}) (map[string]string, error) {
	key := newHashKey(data, id)
	return this.HMGet(key, fieldNames...)
}

//更新多个字段，不指定colNames就更新所有
func (this *RedisCacheClient) OSetFieldsByID(data interface{}, id interface{}, fieldNames ...string) error {
	dataValue := reflect.ValueOf(data).Elem()
	dataType := reflect.TypeOf(data).Elem()
	hash := make(map[interface{}]interface{})
	for _, fieldName := range fieldNames {
		//fieldType, bool := dataType.FieldByName(fieldName)
		field := dataValue.FieldByName(fieldName)
		if !field.IsValid() {
			return errors.New(fmt.Sprintf("[rorm] unexpect field %v-%v err", dataType.Name(), fieldName))
		}
		fieldValue, err := util.GetReflectValue(field)
		if err != nil {
			return errors.New(fmt.Sprintf("[rorm] get field %v-%v:%v err", dataType.Name(), tag, field))
		}

		hash[fieldName] = fieldValue
	}
	if len(hash) > 0 {
		key := newHashKey(data, id)
		return this.HMSet(key, hash)
	} else {
		return errors.New(fmt.Sprintf("[rorm] can not found any field on %v", dataType.Name()))
	}
}

//提取结构体的注解，写入redis
//data 对象
func (this *RedisCacheClient) OSetByID(data interface{}, id interface{}) error {
	dataValue := reflect.ValueOf(data).Elem()
	dataType := reflect.TypeOf(data).Elem()
	hash := make(map[interface{}]interface{})
	for i := 0; i < dataValue.NumField(); i++ {
		field := dataValue.Field(i)
		fieldType := dataType.Field(i)
		tag := fieldType.Tag.Get(tag)
		if tag == ignoreFlag {
			continue
		}
		if tag == "" {
			tag = fieldType.Name
		}
		fieldValue, err := util.GetReflectValue(field)

		if err != nil {
			return errors.New(fmt.Sprintf("[rorm] get field %v-%v:%v err", dataType.Name(), tag, field))
		}
		hash[tag] = fieldValue
	}

	if len(hash) > 0 {
		key := newHashKey(data, id)
		return this.HMSet(key, hash)
	} else {
		return errors.New(fmt.Sprintf("[rorm] can not found any field on %v", dataType.Name()))
	}
}

//获取redis数据，注入结构体
func (this *RedisCacheClient) OGetByID(data interface{}, id interface{}) error {
	dataValue := reflect.ValueOf(data).Elem()
	dataType := reflect.TypeOf(data).Elem()
	key := newHashKey(data, id)
	values, err := this.HGetAll(key)
	if err != nil {
		return err
	}
	for i := 0; i < dataValue.NumField(); i++ {
		field := dataValue.Field(i)
		fieldType := dataType.Field(i)
		tag := fieldType.Tag.Get(tag)
		if tag == ignoreFlag {
			continue
		}
		if tag == "" {
			tag = fieldType.Name
		}
		value := values[tag]
		err := util.SetReflectValue(field, value)
		if err != nil {
			return errors.New(fmt.Sprintf("[rorm] set field %v-%v:%v err", dataType.Name(), tag, value))
		}
	}
	return nil
}
