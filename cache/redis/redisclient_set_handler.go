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

//判断set是否包含成员
func (this *RedisCacheClient) SContains(key string, value interface{}) (bool, error) {
	
	return redis.Bool(this.Do(OP_S_ISMEMBER, key, value))
}

//Set添加数据
func (this *RedisCacheClient) SAddData(key string, value interface{}) error {
	
	_, err := this.Do(OP_S_ADD, key, value)
	return err
}

//随机Set中指定数量的数据   repeat:是否重复
func (this *RedisCacheClient) SRandMember(key string, value int, repeat bool) ([]int, error) {
	
	if repeat {
		value = -value
	}

	return redis.Ints(this.Do(OP_S_RANDMENBER, key, value))
}
