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
	"github.com/garyburd/redigo/redis"
	"os"
	"time"
)

type RedisCacheClient struct {
	MaxIdle     int
	MaxActive   int
	Address     string
	Password    string
	IdleTimeout time.Duration //180 * time.Second
	pool        *redis.Pool
}

//redis.pool.maxActive=200  #最大连接数：能够同时建立的“最大链接个数”

//redis.pool.maxIdle=20     #最大空闲数：空闲链接数大于maxIdle时，将进行回收

//redis.pool.minIdle=5      #最小空闲数：低于minIdle时，将创建新的链接

//redis.pool.maxWait=3000    #最大等待时间：单位ms

func NewRedisClient(config config.CacheConfig) *RedisCacheClient {
	if config.MaxActive == 0 {
		config.MaxActive = 5000
	}
	if config.MaxIdle == 0 {
		config.MaxIdle = 2000
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
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		Address:     config.Address,
		Password:    config.Password,
		IdleTimeout: time.Duration(config.IdleTimeout) * time.Second,
	}
	return redisClient
}

//启动缓存客户端
func (this *RedisCacheClient) Start() {
	this.pool = &redis.Pool{
		MaxIdle:     this.MaxIdle,
		MaxActive:   this.MaxActive,
		IdleTimeout: this.IdleTimeout, //空闲释放时间
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", this.Address)
			if err != nil {
				log.Fatalf("start redis error : %v", err)
				return nil, err
			}
			if this.Password != "" {
				if _, err := c.Do("AUTH", this.Password); err != nil {
					c.Close()
					log.Fatalf("start redis error : %v", err)
					return nil, err
				}
			}
			return c, err
		},
	}
	//测试连接
	err := this.SetData("____test____", "testdata")
	if err != nil {
		log.Fatalf("test redis connection error : %v", err)
	}
}

//关闭缓存客户端
func (this *RedisCacheClient) Close() {
	if this.pool != nil {
		this.pool.Close()
	}
}

func (this *RedisCacheClient) SetNX(key string, value interface{}) (bool, error) {
	conn := this.pool.Get()
	defer conn.Close()
	result, err := redis.Int(conn.Do(OP_SETNX, key, value))
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
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_EXPIRE, key, seconds)
	return err
}

//添加数据
func (this *RedisCacheClient) SetData(key string, value interface{}) error {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_SET, key, value)
	return err
}

func (this *RedisCacheClient) Incr(key string) (int, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do(OP_INCR, key))
}

func (this *RedisCacheClient) Decr(key string) (int, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do(OP_DECR, key))
}

func (this *RedisCacheClient) SelectDB(dbNumber int) error {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_SELECT, dbNumber)
	return err
}

func (this *RedisCacheClient) GetDataInt32(key string) (int32, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return Int32(redis.Int(conn.Do(OP_GET, key)))
}

func (this *RedisCacheClient) GetDataInt64(key string) (int64, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do(OP_GET, key))
}

//获取数据
func (this *RedisCacheClient) GetData(key string) (string, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.String(conn.Do(OP_GET, key))
}

//导出数据
func (this *RedisCacheClient) Dump(key string) (string, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.String(conn.Do(OP_DUMP, key))
}

//导入数据
func (this *RedisCacheClient) Restore(key string, data string) (string, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.String(conn.Do(OP_RESTORE, key, 0, data))
}

//是否存在数据
func (this *RedisCacheClient) Exists(key string) (bool, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do(OP_EXISTS, key))
}

//添加数据
func (this *RedisCacheClient) DelData(key string) error {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_DEL, key)
	return err
}

//清除所有数据
func (this *RedisCacheClient) FlashAll() error {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_FLUSHALL)
	return err
}

//// 存map
//func (this *RedisCacheClient)SetMap(key string ,value map[string]string) bool{
//	conn := this.pool.Get()
//	defer conn.Close()
//	// 转换成json
//	v,_ := json.Marshal(value)
//	// 存redis
//	_,err := conn.Do("SETNX",key, v)
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
//	value,err := redis.Bytes(conn.Do("GET",key))
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
func (this *RedisCacheClient) Subscribe(callback func(channel, value string), channel ...interface{}) {
	//defer conn.Close()
	psc := redis.PubSubConn{Conn: this.pool.Get()}
	psc.Subscribe(channel...)
	go func() {
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
	}()
}

func (this *RedisCacheClient) PSubscribe(callback func(pattern, channel, value string), channel ...interface{}) {
	//defer conn.Close()
	psc := &redis.PubSubConn{Conn: this.pool.Get()}
	psc.PSubscribe(channel...)
	go func() {
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
	}()
}

//发布数据
func (this *RedisCacheClient) Publish(channel, value interface{}) {
	conn := this.pool.Get()
	defer conn.Close()
	conn.Do(OP_PUBLISH, channel, value)
}
