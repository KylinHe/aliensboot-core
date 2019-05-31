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
	"errors"
)

type Rank struct {
	Member string
	Score  int64
}

func Ranks(result interface{}, err error) ([]Rank, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: IntMap expects even number of values result")
	}
	ranks := []Rank{}
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].([]byte)
		if !ok {
			return nil, errors.New("redigo: ScanMap key not a bulk string value")
		}
		value, err := redis.Int(values[i+1], nil)
		if err != nil {
			return nil, err
		}
		ranks = append(ranks, Rank{Score: int64(value), Member: string(key)})
	}
	return ranks, nil
}

//ZADD key score member [[score member] [score member] ...]
//将一个或多个 member 元素及其 score 值加入到有序集 key 当中。
func (this *RedisCacheClient) ZAdd(key string, score int64, member interface{}) error {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_Z_ADD, key, score, member)
	return err
}

//ZREM key member [member ...]
//
//移除有序集 key 中的一个或多个成员，不存在的成员将被忽略。
//
//当 key 存在但不是有序集类型时，返回一个错误。
func (this *RedisCacheClient) ZRem(key string, member interface{}) error {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_Z_REM, key, member)
	return err
}

func (this *RedisCacheClient) ZCount(key string, min int32, max int32) (int, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do(OP_Z_COUNT, key, min, max))
}


////ZADD key score member [[score member] [score member] ...]
////添加多条有序集合
//func (this *RedisCacheClient)ZAdds(key string, scoreMembers ...string) bool{
//	conn := this.pool.Get()
//	defer conn.Close()
//	_ , err := conn.Do(OP_Z_ADD, key, scoreMembers...)
//	if  err != nil{
//		//log.Debug("%v",err)
//		return false
//	}
//	return true
//}

/**
ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]

返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列。

具有相同 score 值的成员按字典序(lexicographical order)来排列(该属性是有序集提供的，不需要额外的计算)。

可选的 LIMIT 参数指定返回结果的数量及区间(就像SQL中的 SELECT LIMIT offset, count )，注意当 offset 很大时，定位 offset 的操作可能需要遍历整个有序集，此过程最坏复杂度为 O(N) 时间。

可选的 WITHSCORES 参数决定结果集是单单返回有序集的成员，还是将有序集成员及其 score 值一起返回。
该选项自 Redis 2.0 版本起可用。
*/
func (this *RedisCacheClient) ZRangeByScoreLimit(key string, min int32, max int32, offset int32, count int32) ([]Rank, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return Ranks(conn.Do(OP_Z_RANGEBYSCORE, key, min, max, PARAM_WITHSCORES, PARAM_LIMIT, offset, count))
}

//获取score小于 {max}的{count}个信息
func (this *RedisCacheClient) ZRevRangeByScoreBeforeLimit(key string, max int32, count int32) ([]string, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do(OP_Z_REVRANGEBYSCORE, key, max, PARAM_Z_MIN, PARAM_LIMIT, 0, count))
}

//获取score小于 {max}的{count}个信息
//func (this *RedisCacheClient)ZRangeByScoreBeforeLimit(key string, max int32, count int32) []string {
//	conn := this.pool.Get()
//	defer conn.Close()
//	result,err := redis.Strings(conn.Do(OP_Z_RANGEBYSCORE,key,PARAM_Z_MAX,max,PARAM_LIMIT,0,count))
//	if  err != nil{
//		//log.Debug("%v",err)
//	}
//	return result
//}

//ZRANGE key start stop [WITHSCORES]   start:0  stop :-1 显示所有 下标从0开始
//返回有序集 key 中，指定区间内的成员。
//其中成员的位置按 score 值递增(从小到大)来排序。
func (this *RedisCacheClient) ZRange(key string, start int32, stop int32) ([]string, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do(OP_Z_RANGE, key, start, stop))
}

//ZREVRANGE key start stop [WITHSCORES]
//返回有序集 key 中，指定区间内的成员。
//其中成员的位置按 score 值递减(从大到小)来排列
func (this *RedisCacheClient) ZRevRange(key string, start int32, stop int32) ([]string, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do(OP_Z_REVRANGE, key, start, stop))
}

func (this *RedisCacheClient) ZRangeWithScore(key string, start int32, stop int32) ([]Rank, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return Ranks(conn.Do(OP_Z_RANGE, key, start, stop, PARAM_WITHSCORES))
}

func (this *RedisCacheClient) ZRevRangeWithScore(key string, start int32, stop int32) ([]Rank, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return Ranks(conn.Do(OP_Z_REVRANGE, key, start, stop, PARAM_WITHSCORES))
}

func (this *RedisCacheClient) ZRevRank(key string, member interface{}) (int, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do(OP_Z_REVRANK, key, member))
}
