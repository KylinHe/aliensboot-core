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
	//"github.com/name5566/leaf/log"
)

//获取列表长度
func (this *RedisCacheClient) LLen(key string) (int, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do(OP_L_LEN, key))
}

//列表添加多条数据
func (this *RedisCacheClient) LPushString(key string, value string) (int, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do(OP_L_PUSH, key, value))
}

func (this *RedisCacheClient) LPush(key string, value interface{}) (int, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do(OP_L_PUSH, key, value))
}

//func (this *RedisCacheClient)LPushMulti(key string, values []string) int {
//	conn := this.pool.Get()
//	defer conn.Close()
//	//for i := 0; i < len(values) ; i++  {
//	//	conn.Do(OP_L_PUSH,key,value...)
//	//}
//
//	len, err := redis.Int(conn.Do(OP_L_PUSH, key, values...))
//	if  err != nil{
//		//log.Debug("%v",err)
//	}
//	return len
//}

//func (this *RedisCacheClient)RPush(key string,value ...string) int {
//	conn := this.pool.Get()
//	defer conn.Close()
//	//for i := 0; i < len(values) ; i++  {
//	//	conn.Do(OP_L_PUSH,key,value...)
//	//}
//	len, err := redis.Int(conn.Do(OP_R_PUSH,key,value...))
//	if  err != nil{
//		//log.Debug("%v",err)
//	}
//	return len
//}

//向列条添加一条数据
//func (this *RedisCacheClient)LPush(key string,value string) {
//	conn := this.pool.Get()
//	defer conn.Close()
//	conn.Do(OP_L_PUSH,key,value)
//}

//获取列表所有数据
func (this *RedisCacheClient) LRangeAll(key string) ([]string, error) {
	return this.LRange(key, 0, -1)
}

func (this *RedisCacheClient) LRangeAllByte(key string) ([][]byte, error) {
	return this.LRangeBytes(key, 0, -1)
}

//获取列表指定范围内的数据
func (this *RedisCacheClient) LRange(key string, star int, stop int) ([]string, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do(OP_L_RANGE, key, star, stop))
}

func (this *RedisCacheClient) LRangeBytes(key string, star int, stop int) ([][]byte, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.ByteSlices(conn.Do(OP_L_RANGE, key, star, stop))
}
