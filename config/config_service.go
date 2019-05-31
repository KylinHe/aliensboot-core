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

type ServiceConfig struct {
	ID       string //服务器的id
	Name     string //服务名称
	Address  string //服务地址 域名或ip
	Port     int    //服务端端口
	Unique   bool   //是否全局唯一
	Protocol string //提供服务的协议 GRPC HTTP WBSOCKET
	Lbs      string //负载均衡策略
	Local    bool   //是否支持本地调用 调用方和服务方在一个进程。优化作用
}
