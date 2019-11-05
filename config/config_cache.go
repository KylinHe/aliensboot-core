/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/10/25
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package config

type CacheConfig struct {
	Address     string
	Password    string
	MaxActive   int      //最大连接数，即最多的tcp连接数，一般建议往大的配置，但不要超过操作系统文件句柄个数（centos下可以ulimit -n查看）
	MaxIdle     int      //最大空闲连接数，即会有这么多个连接提前等待着，但过了超时时间也会关闭
	IdleTimeout int      //空闲连接超时时间，但应该设置比redis服务器超时时间短。否则服务端超时了，客户端保持着连接也没用
	Wait        bool     //如果超过最大连接，是报错，还是等待
}
