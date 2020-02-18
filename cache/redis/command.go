/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/8/15
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package redis

type Command struct {
	Args []interface{}
}

//操作类型
const (
	TTL_VALUE_FOREVER int = -1
	TTL_VALUE_NOTFOUND int = -2

	PARAM_WITHSCORES string = "WITHSCORES"
	PARAM_LIMIT      string = "limit"
	PARAM_Z_MAX      string = "+inf"
	PARAM_Z_MIN      string = "-inf"

	OP_SELECT string = "SELECT"
	OP_PING string = "PING"
	OP_AUTH string = "AUTH"
	OP_MULTI string = "MULTI"
	OP_EXEC string = "EXEC"

	OP_S_ADD        string = "SADD"
	OP_S_RANDMENBER string = "SRANDMEMBER"
	OP_S_ISMEMBER   string = "SISMEMBER"

	OP_DUMP    string = "DUMP"
	OP_RESTORE string = "RESTORE"

	OP_SET    string = "SET"
	OP_GET    string = "GET"
	OP_DEL    string = "DEL"
	OP_EXISTS string = "EXISTS"
	OP_SETEX  string = "SETEX"
	OP_SETNX  string = "SETNX"
	OP_INCR   string = "INCR" //自增长
	OP_DECR   string = "DECR" //自减

	OP_EXPIRE string = "EXPIRE"
	OP_TTL string = "TTL"
	OP_SCAN    string = "SCAN"


	OP_FLUSHALL string = "FLUSHALL"

	OP_H_SET     string = "HSET"
	OP_H_GET     string = "HGET"
	OP_H_GETALL  string = "HGETALL"
	OP_H_MGET    string = "HMGET"
	OP_H_MSET    string = "HMSET"
	OP_H_DEL     string = "HDEL"
	OP_H_HINCRBY string = "HINCRBY"
	OP_H_EXISTS  string = "HEXISTS"
	OP_H_LEN     string = "HLEN"

	OP_L_PUSH  string = "LPUSH"
	OP_R_PUSH  string = "RPUSH"
	OP_L_RANGE string = "LRANGE"
	OP_L_LEN   string = "LLEN"
	OP_L_REM   string = "LREM"

	OP_Z_COUNT           string = "ZCOUNT"
	OP_Z_ADD             string = "ZADD"
	OP_Z_REM             string = "ZREM"
	OP_Z_RANGE           string = "ZRANGE"
	OP_Z_RANGEBYSCORE    string = "ZRANGEBYSCORE"
	OP_Z_REVRANGEBYSCORE string = "ZREVRANGEBYSCORE"
	OP_Z_REVRANGE        string = "ZREVRANGE"

	OP_Z_REVRANK string = "ZREVRANK"

	OP_PUBLISH string = "PUBLISH"
)
