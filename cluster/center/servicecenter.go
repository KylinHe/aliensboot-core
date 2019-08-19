/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/6/4
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package center

import (
	"github.com/KylinHe/aliensboot-core/cluster/center/service"
	"github.com/KylinHe/aliensboot-core/common/util"
	"github.com/KylinHe/aliensboot-core/config"
	"github.com/KylinHe/aliensboot-core/log"
	"os"
	"strconv"
)

var ClusterCenter ServiceCenter = &ETCDServiceCenter{} //服务调度中心

const NodeSplit string = "/"

const ServiceNodeName string = "service"

const ConfigNodeName string = "config"

type DataEventType int32

const (
	PUT    DataEventType = 0
	DELETE DataEventType = 1
)

func PublicService(config config.ServiceConfig, handler interface{}) service.IService {
	var service service.IService = nil
	service = startService(config, handler)
	if ClusterCenter.IsConnect() {
		//RPC启动成功,则发布到中心服务器
		if !ClusterCenter.PublicService(service, config) {
			log.Fatalf("service %v can not be public", service.GetName())
		}
	} else {
		log.Infof(config.Name + " cluster center is not connected")
	}
	return service
}

func startService(config config.ServiceConfig, handler interface{}) service.IService {
	config.ID = ClusterCenter.GetNodeID() //节点id

	rpcAddress := os.Getenv("ServiceAddress")
	if rpcAddress != "" {
		config.Address = rpcAddress
	}

	rpcPort := os.Getenv("ServicePort")
	if rpcPort != "" {
		newPort, err := strconv.Atoi(rpcPort)
		if err == nil {
			config.Port	= newPort
		}
	}

	//地址没有发布到外网 采用内网地址
	if config.Address == "" {
		config.Address = util.GetIP()
	}
	service, err := service.NewService(config)
	if err != nil {
		log.Fatalf("create service err : %v", err)
	}
	service.SetHandler(handler)
	if !service.Start() {
		log.Fatalf("service %v can not be start", service.GetName())
	}
	return service
}

func ReleaseService(service service.IService) {
	//if !ClusterCenter.IsConnect() {
	//	log.Errorf(" cluster center is not connected")
	//	return
	//}
	//先从中心释放，再内部关闭，缓解关闭期间其他服务请求转发过来
	ClusterCenter.ReleaseService(service)
	if service != nil {
		service.Close()
	}
}

type ConfigListener func(data []byte)

//
type DataPrefixListener func(eventType DataEventType, data []byte, dataRootName string, dataName string)


type ServiceCenter interface {

	GetNodeID() string //获取当前节点id

	ConnectCluster(config config.ClusterConfig)

	Close()

	IsConnect() bool

	PublicConfig(configName string, content []byte) bool        //发布配置

	SubscribeConfig(configName string, listener ConfigListener) //订阅配置

	SubscribeConfigWithPrefix(configName string, listener DataPrefixListener) //

	ReleaseService(service service.IService)                  //释放服务

	PublicService(service service.IService, serviceConfig config.ServiceConfig) bool //发布服务

	SubscribeServices(serviceName ...string)  //订阅服务

	GetAllService(serviceName string) []service.IService //获取所有的服务

	GetService(serviceName string, serviceID string) service.IService //获取指定服务

	AllocService(serviceName string, param string) service.IService   //按照负载均衡策略 分配一个可用的服务


}
