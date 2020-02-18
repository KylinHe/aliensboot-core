/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/3/27
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package redis

import (
	"github.com/KylinHe/aliensboot-core/config"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/KylinHe/aliensboot-core/task"
	"github.com/garyburd/redigo/redis"
	"os"
	"time"
)

type ErrorHandler func(err error, command string, args ...interface{})

var ErrNil = redis.ErrNil

type RedisCacheClient struct {
	config.CacheConfig
	IdleTimeout time.Duration //180 * time.Second
	pool        *redis.Pool
	errorHandler ErrorHandler
}


//redis.pool.maxActive=200  #最大连接数：能够同时建立的“最大链接个数”

//redis.pool.maxIdle=20     #最大空闲数：空闲链接数大于maxIdle时，将进行回收

//redis.pool.minIdle=5      #最小空闲数：低于minIdle时，将创建新的链接

//redis.pool.maxWait=3000    #最大等待时间：单位ms

func NewRedisClient(config config.CacheConfig) *RedisCacheClient {
	if config.MaxActive == 0 {
		config.MaxActive = 2000
	}
	if config.MaxIdle == 0 {
		config.MaxIdle = 500
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = 120
	}

	//优先使用环境变量
	address := os.Getenv("CacheAddress")
	password := os.Getenv("CachePassword")
	if address != "" {
		config.Address = address
		config.Password = password
	}

	redisClient := &RedisCacheClient{
		CacheConfig: config,
		IdleTimeout: time.Duration(config.IdleTimeout) * time.Second,
	}
	return redisClient
}

//启动缓存客户端
func (this *RedisCacheClient) Start() {
	this.pool = &redis.Pool{
		MaxIdle:     this.MaxIdle,
		MaxActive:   this.MaxActive,
		Wait: this.Wait,
		IdleTimeout: this.IdleTimeout, //空闲释放时间
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", this.Address)
			if err != nil {
				return nil, err
			}
			if this.Password != "" {
				if _, err := c.Do(OP_AUTH, this.Password); err != nil {
					return nil, err
				}
			}
			return c, err
		},
	}
	//测试连接
	_, err := this.Do(OP_PING)
	if err != nil {
		log.Fatalf("test redis connection error : %v", err)
	}
}

//关闭缓存客户端
func (this *RedisCacheClient) Close() error {
	if this.pool != nil {
		return this.pool.Close()
	}
	return nil
}

func (this *RedisCacheClient) SetErrorHandler(handler ErrorHandler)  {
	this.errorHandler = handler
}


// 管道执行命令
func (this *RedisCacheClient) PipelineCommands(commands []Command) error {
	conn := this.pool.Get()
	defer conn.Close()
	for _, cmd := range commands {
		if len(cmd.Args) == 0 {
			continue
		}
		if err := conn.Send(cmd.Args[0].(string), cmd.Args[1:]...); err != nil {
			this.errorHandler(err, cmd.Args[0].(string), cmd.Args[1:]...)
			return err
		}
	}
	if err := conn.Flush(); err != nil {
		this.errorHandler(err, "Flush")
		return err
	}
	return nil
}

// 事务执行命令
func (this *RedisCacheClient) MultiCommands(commands []Command) error {
	conn := this.pool.Get()
	defer conn.Close()
	err := conn.Send(OP_MULTI)
	if err != nil {
		this.errorHandler(err, OP_MULTI)
		return err
	}
	for _, cmd := range commands {
		if len(cmd.Args) == 0 {
			continue
		}
		if err := conn.Send(cmd.Args[0].(string), cmd.Args[1:]...); err != nil {
			this.errorHandler(err, cmd.Args[0].(string), cmd.Args[1:]...)
			return err
		}
	}
	err = conn.Send(OP_EXEC)
	if err != nil {
		this.errorHandler(err, OP_EXEC)
	}
	return err
}

func (this *RedisCacheClient) Do(command string, args ...interface{}) (reply interface{}, err error) {
	conn := this.pool.Get()
	defer conn.Close()
	replay, err := conn.Do(command, args...)
	if err != nil && this.errorHandler != nil {
		this.errorHandler(err, command, args...)
	}
	return replay, err
}

func (this *RedisCacheClient) SetNX(key string, value interface{}) (bool, error) {
	result, err := redis.Int(this.Do(OP_SETNX, key, value))
	if err != nil {
		return false, err
	}
	return result == 1, err

}

func Int32(param int, err error) (int32, error) {
	return int32(param), err
}

//设置数据过期时间
func (this *RedisCacheClient) Expire(key string, seconds int) error {
	
	_, err := this.Do(OP_EXPIRE, key, seconds)
	return err
}

func (this *RedisCacheClient) Ttl(key string) (int, error) {

	return redis.Int(this.Do(OP_TTL, key))
}


//添加数据
func (this *RedisCacheClient) SetData(key string, value interface{}) error {
	
	_, err := this.Do(OP_SET, key, value)
	return err
}

func (this *RedisCacheClient) Incr(key string) (int, error) {
	
	return redis.Int(this.Do(OP_INCR, key))
}

func (this *RedisCacheClient) Decr(key string) (int, error) {
	
	return redis.Int(this.Do(OP_DECR, key))
}

func (this *RedisCacheClient) Scan(cursor int, match interface{},count int) (int, []string, error) {
	var result []string
	value, err := redis.Values(this.Do(OP_SCAN, cursor, "MATCH", match, "COUNT",count))
	if err != nil {
		return cursor,result,err
	}
	for i, v := range value {
		if i == 0 {	// idx0 是当前游标
			cursor,err = redis.Int(v,err)
			if err != nil {
				continue
			}
		} else {	// idx1 是数值,需要获取的结果
			values,err := redis.Strings(v,err)
			if err != nil {
				continue
			}
			result = append(result,values...)
		}
	}
	return cursor,result, nil
}

func (this *RedisCacheClient) SelectDB(dbNumber int) error {
	
	_, err := this.Do(OP_SELECT, dbNumber)
	return err
}

func (this *RedisCacheClient) GetDataInt32(key string) (int32, error) {
	
	return Int32(redis.Int(this.Do(OP_GET, key)))
}

func (this *RedisCacheClient) GetDataInt64(key string) (int64, error) {
	
	return redis.Int64(this.Do(OP_GET, key))
}

//获取数据
func (this *RedisCacheClient) GetData(key string) (string, error) {
	
	return redis.String(this.Do(OP_GET, key))
}

//导出数据
func (this *RedisCacheClient) Dump(key string) (string, error) {
	
	return redis.String(this.Do(OP_DUMP, key))
}

//导入数据
func (this *RedisCacheClient) Restore(key string, data string) (string, error) {
	
	return redis.String(this.Do(OP_RESTORE, key, 0, data))
}

//是否存在数据
func (this *RedisCacheClient) Exists(key string) (bool, error) {
	
	return redis.Bool(this.Do(OP_EXISTS, key))
}

//添加数据
func (this *RedisCacheClient) DelData(key string) error {
	
	_, err := this.Do(OP_DEL, key)
	return err
}

//清除所有数据
func (this *RedisCacheClient) FlashAll() error {
	
	_, err := this.Do(OP_FLUSHALL)
	return err
}

//// 存map
//func (this *RedisCacheClient)SetMap(key string ,value map[string]string) bool{
//	conn := this.pool.Get()
//	defer conn.Close()
//	// 转换成json
//	v,_ := json.Marshal(value)
//	// 存redis
//	_,err := this.Do("SETNX",key, v)
//	if err != nil {
//		//log.Debug("%v",err)
//		return false
//	}
//	return true
//}
//
//// 取map
//func (this *RedisCacheClient)GetMap(key string) map[string]string {
//	conn := this.pool.Get()
//	defer conn.Close()
//	var imap map[string]string
//	value,err := redis.Bytes(this.Do("GET",key))
//	if err != nil {
//		//log.Debug("%v",err)
//		return nil
//	}
//	// json转map
//	errShal := json.Unmarshal(value,&imap)
//	if errShal != nil {
//		//log.Debug("%v",errShal)
//		return nil
//	}
//	return imap
//}

//订阅数据变更
func (this *RedisCacheClient) Subscribe(callback func(channel, value string), channel ...interface{}) error {
	//defer conn.Close()
	psc := redis.PubSubConn{Conn: this.pool.Get()}
	err := psc.Subscribe(channel...)
	task.SafeGo(func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				value, _ := redis.String(v.Data, nil)
				callback(v.Channel, value)
			case error:
				//log.Debug("error: %v\n", v)
				return
			}
		}
	})

	return err
}

func (this *RedisCacheClient) PSubscribe(callback func(pattern, channel, value string), channel ...interface{}) error {
	//defer conn.Close()
	psc := &redis.PubSubConn{Conn: this.pool.Get()}
	err := psc.PSubscribe(channel...)
	task.SafeGo(func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.PMessage:
				value, _ := redis.String(v.Data, nil)
				callback(v.Pattern, v.Channel, value)
			case error:
				//log.Debug("error: %v\n", v)
				return
			}
		}
	})
	return err
}

//发布数据
func (this *RedisCacheClient) Publish(channel, value interface{}) error {
	_, err := this.Do(OP_PUBLISH, channel, value)
	return err
}
