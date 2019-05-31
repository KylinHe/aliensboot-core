/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/3/24
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package service

import "github.com/KylinHe/aliensboot-core/protocol/base"

//调用方式
const (
	WEBSOCKET string = "websocket"
	HTTP      string = "http"
	GRPC      string = "grpc"
)

//type Callback func(response *base.Any, err error)

type IService interface {
	GetID() string
	SetID(id string)

	SetName(name string)
	GetName() string

	GetAddress() string
	GetPort() int
	//GetConfig() *CenterService
	Start() bool   //启动服务
	Connect() bool //连接服务
	Close()        //关闭服务

	Equals(other IService) bool //比较服务
	IsLocal() bool              //是否本机服务

	Request(request *base.Any) (*base.Any, error) //同步请求 阻塞
	//AsyncRequest(asyncCall *AsyncCall) //异步请求，响应采用回调
	Send(request *base.Any) error //发送消息，不需要响应

	SetHandler(handler interface{}) //设置处理句柄

	//GetLBS() string
	//KickOut(request interface{}) error //服务推送
}

type CenterService struct {
	Address  string `json:"address"`  //服务地址 域名或ip
	Port     int    `json:"port"`     //服务端端口
	Protocol string `json:"protocol"` //提供服务的协议 GRPC HTTP WBSOCKET
	ID       string `json:"-"`        //服务器的id
	Name     string `json:"-"`        //服务名称
	Unique   bool   `json:"-"`        //是否全局唯一
}

func (this *CenterService) GetAddress() string {
	return this.Address
}

func (this *CenterService) SetAddress(address string) {
	this.Address = address
}

func (this *CenterService) GetPort() int {
	return this.Port
}

func (this *CenterService) SetPort(port int) {
	this.Port = port
}

func (this *CenterService) GetID() string {
	return this.ID
}

func (this *CenterService) SetID(id string) {
	this.ID = id
}

func (this *CenterService) GetName() string {
	return this.Name
}

func (this *CenterService) SetName(name string) {
	this.Name = name
}

//
//type centerService struct {
//	Ip       string `json:"ip"`
//	Port     int	`json:"port"`
//	//Address  string `json:"address"`  	//服务访问地址 写入到中心服务器供外部调用
//	Protocol string `json:"protocol"` 		//服务的访问方式
//	id       string `json:"-"`              //服务ID
//	name     string `json:"-"`              //服务类型
//}
