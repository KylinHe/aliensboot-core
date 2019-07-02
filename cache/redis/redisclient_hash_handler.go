/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/3/29
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package redis

import (
	"github.com/garyburd/redigo/redis"
)

//自增指定数量
func (this *RedisCacheClient) HIncrby(key interface{}, field string, increment int) (int, error) {
	
	return redis.Int(this.Do(OP_H_HINCRBY, key, field, increment))
}

func (this *RedisCacheClient) HSet(key interface{}, field string, value interface{}) error {
	
	_, err := this.Do(OP_H_SET, key, field, value)
	return err
}

func (this *RedisCacheClient) HGetBytes(key interface{}, field string) ([]byte, error) {
	
	return redis.Bytes(this.Do(OP_H_GET, key, field))
}

func (this *RedisCacheClient) HGet(key interface{}, field string) (string, error) {
	
	return redis.String(this.Do(OP_H_GET, key, field))
}

func (this *RedisCacheClient) HDel(key interface{}, field string) error {
	
	_, err := this.Do(OP_H_DEL, key, field)
	return err
}

//func (this *RedisCacheClient) HSetData(key interface{}, data interface{}) {
//	this.HSetFieldData(key, "", data)
//}
//
////获取redis数据，注入结构体
//func (this *RedisCacheClient) HGetData(key interface{}, data interface{}) {
//	this.HGetFieldData(key, "", data)
//}

func (this *RedisCacheClient) HGetBool(key interface{}, field string) (bool, error) {
	
	return redis.Bool(this.Do(OP_H_GET, key, field))
}

func (this *RedisCacheClient) HGetInt32(key interface{}, field string) (int32, error) {
	
	return Int32(redis.Int(this.Do(OP_H_GET, key, field)))
}

func (this *RedisCacheClient) HGetInt64(key interface{}, field string) (int64, error) {
	
	return redis.Int64(this.Do(OP_H_GET, key, field))
}

//判断hash字段是否存在
func (this *RedisCacheClient) HFieldExists(key interface{}, field string) (bool, error) {
	
	return redis.Bool(this.Do(OP_H_EXISTS, key, field))
}

//获取所有的hash数据
func (this *RedisCacheClient) HGetAllInt(key interface{}) (map[string]int, error) {
	
	return redis.IntMap(this.Do(OP_H_GETALL, key))
}

func (this *RedisCacheClient) HGetAllInt64(key interface{}) (map[string]int64, error) {
	
	return redis.Int64Map(this.Do(OP_H_GETALL, key))
}

func (this *RedisCacheClient) HGetAll(key interface{}) (map[string]string, error) {
	
	return redis.StringMap(this.Do(OP_H_GETALL, key))
}

//批量添加字段
func (this *RedisCacheClient) HMSet(key interface{}, fields map[interface{}]interface{}) error {
	

	params := []interface{}{key}
	for key, value := range fields {
		params = append(params, key, value)
	}
	_, err := this.Do(OP_H_MSET, params...)
	return err
}

func (this *RedisCacheClient) HMGet(key interface{}, fieldNames ...interface{}) (map[string]string, error) {
	

	params := []interface{}{key}
	params = append(params, fieldNames...)

	values, err := redis.Strings(this.Do(OP_H_MGET, params...))
	if err != nil {
		return nil, err
	}
	//if len(values) != len(fieldNames) {
	//	return nil, errors
	//}
	results := make(map[string]string)
	for index, value := range values {
		results[fieldNames[index].(string)] = value
	}
	return results, nil
}
