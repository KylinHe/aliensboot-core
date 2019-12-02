/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 *
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package center

//服务中心，处理服务的调度和查询
//import (
//	"encoding/json"
//	"github.com/KylinHe/aliensboot-core/cluster/center/service"
//	"github.com/KylinHe/aliensboot-core/config"
//	"github.com/KylinHe/aliensboot-core/log"
//	"github.com/KylinHe/aliensboot-core/task"
//	"github.com/samuel/go-zookeeper/zk"
//	"gopkg.in/mgo.v2/bson"
//	"time"
//)
//
//type ZKServiceCenter struct {
//	*service.Container
//	zkCon       *zk.Conn
//	zkName      string
//	serviceRoot string
//	configRoot  string
//
//	nodeId string //当前集群节点的id
//	//lbs string //default polling
//	certFile   string
//	keyFile    string
//	commonName string
//}
//
//func (this *ZKServiceCenter) GetNodeID() string {
//	return this.nodeId
//}
//
////启动服务中心客户端
////func (this *ZKServiceCenter) Connect(address string, timeout int, zkName string, nodeID string) {
////	this.ConnectCluster([]string{address}, timeout, zkName, nodeID)
////}
//
//func (this *ZKServiceCenter) ConnectCluster(config config.ClusterConfig) {
//	if config.ID == "" {
//		config.ID = bson.NewObjectId().Hex()
//		//config.ID =
//		//panic("cluster nodeID can not be empty")
//	}
//	if config.Timeout == 0 {
//		config.Timeout = 10
//	}
//	//this.lbs = config.LBS
//	this.zkName = config.Name
//	this.nodeId = config.ID
//	//this.certFile = config.CertFile
//	//this.keyFile = config.KeyFile
//	//this.commonName = config.CommonName
//	//this.serviceFactory = serviceFactory
//	c, _, err := zk.Connect(config.Servers, time.Duration(config.Timeout)*time.Second)
//	if err != nil {
//		panic(err)
//	}
//	this.Container = service.NewContainer()
//	this.serviceRoot = NodeSplit + this.zkName + NodeSplit + ServiceNodeName
//	this.configRoot = NodeSplit + this.zkName + NodeSplit + ConfigNodeName
//
//	this.zkCon = c
//	this.confirmNode(NodeSplit + this.zkName)
//	this.confirmNode(this.serviceRoot)
//}
//
//func (this *ZKServiceCenter) IsConnect() bool {
//	return this.zkCon != nil
//}
//
//func (this *ZKServiceCenter) assert() {
//	if this.zkCon == nil {
//		panic("mast start service center first")
//	}
//}
//
////关闭服务中心
//func (this *ZKServiceCenter) Close() {
//	if this.zkCon != nil {
//		this.zkCon.Close()
//	}
//}
//
////订阅服务  能实时更新服务信息
//func (this *ZKServiceCenter) SubscribeServices(serviceTypes ...string) {
//	this.assert()
//	for _, serviceType := range serviceTypes {
//		this.SubscribeService(serviceType)
//	}
//}
//
//func (this *ZKServiceCenter) ReleaseService(service service.IService) {
//
//}
//
//func (this *ZKServiceCenter) SubscribeService(serviceName string) {
//	//this.SubscribeConfig("lbs"+NodeSplit+serviceName, func(data []byte) {
//	//	this.Container.SetLbs(serviceName, string(data))
//	//})
//	path := this.serviceRoot + NodeSplit + serviceName
//	//desc := this.confirmContentNode(path)
//	serviceIDs, _, ch, err := this.zkCon.ChildrenW(path)
//	if err != nil {
//		log.Errorf("subscribe service %v error: %v", path, err)
//		return
//	}
//	services := []service.IService{}
//	for _, serviceID := range serviceIDs {
//		servicePath := path + NodeSplit + serviceID
//		data, _, err := this.zkCon.Get(servicePath)
//		if err != nil {
//			log.Errorf("get service %v data error: %v", servicePath, err)
//			continue
//		}
//		centerService := &service.CenterService{}
//		err1 := json.Unmarshal(data, centerService)
//		if err1 != nil {
//			log.Errorf("unmarshal service %v data error: %v", servicePath, err1)
//			continue
//		}
//		service, _ := service.NewService2(centerService, serviceID, serviceName)
//		services = append(services, service)
//	}
//
//	this.UpdateServices(serviceName, services)
//	go this.openListener(serviceName, path, ch)
//}
//
//func (this *ZKServiceCenter) openListener(serviceType string, path string, ch <-chan zk.Event) {
//	event, _ := <-ch
//	//更新服务节点信息
//	if event.Type == zk.EventNodeChildrenChanged {
//		this.SubscribeService(serviceType)
//	}
//}
//
////
//func (this *ZKServiceCenter) confirmNode(path string, flags ...int32) bool {
//	_, err := this.zkCon.Create(path, nil, 0, zk.WorldACL(zk.PermAll))
//	return err == nil
//}
//
//func (this *ZKServiceCenter) confirmContentNode(path string, flags ...int32) string {
//	_, err := this.zkCon.Create(path, nil, 0, zk.WorldACL(zk.PermAll))
//	if err != nil {
//		data, _, _ := this.zkCon.Get(path)
//		return string(data)
//	}
//	return ""
//}
//
//func (this *ZKServiceCenter) confirmDataNode(path string, data []byte) bool {
//	byteData := []byte(data)
//	_, err := this.zkCon.Create(path, byteData, 0, zk.WorldACL(zk.PermAll))
//	if err != nil {
//		this.zkCon.Set(path, byteData, -1)
//	}
//	return err == nil
//}
//
////发布服务
//func (this *ZKServiceCenter) PublicService(service service.IService, config config.ServiceConfig) bool {
//	this.assert()
//	if !service.IsLocal() {
//		log.Error("service info is invalid")
//		return false
//	}
//	//path string, data []byte, version int32
//	data, err := json.Marshal(service)
//	if err != nil {
//		log.Errorf("marshal json service data error : %v", err)
//		return false
//	}
//	serviceName := service.GetName()
//	serviceId := service.GetID()
//	servicePath := this.serviceRoot + NodeSplit + serviceName
//	if config.Unique {
//		//TODO 可能有事务上的问题 需要优化
//		child, _, _ := this.zkCon.Children(servicePath)
//		if child != nil && len(child) > 0 {
//			log.Errorf("unique service %v-%v already exist.", serviceName, child)
//			return false
//		}
//	}
//
//	this.confirmNode(servicePath)
//	id, err := this.zkCon.Create(servicePath+NodeSplit+serviceId, data,
//		zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
//	if err != nil {
//		log.Errorf("create service error : %v", err)
//		return false
//	}
//	log.Infof("public %v success : %v", serviceName, id)
//	//服务注册在容器
//	this.UpdateService(service, true)
//	return true
//}
//
////发布配置信息
//func (this *ZKServiceCenter) PublicConfig(configType string, configContent []byte) bool {
//	this.assert()
//	if configType == "" {
//		log.Info("config type con not be empty")
//		return false
//	}
//	configPath := this.configRoot + NodeSplit + configType
//	this.confirmNode(configPath)
//	_, err := this.zkCon.Set(configPath, configContent, -1)
//	if err != nil {
//		log.Info("public config %v  err : %v", configType, err)
//		return false
//	}
//	log.Info("public config %v success", configType)
//	return true
//}
//
////订阅服务  能实时更新服务信息
//func (this *ZKServiceCenter) SubscribeConfig(configName string, configHandler DataListener) {
//	this.assert()
//	path := this.configRoot + NodeSplit + configName
//	this.confirmNode(path)
//	content, _, ch, err := this.zkCon.GetW(path)
//	if err != nil {
//		log.Info("subscribe config %v error: %v", path, err)
//		return
//	}
//	configHandler(content)
//	task.SafeGo(func() {
//		for {
//			event, _ := <-ch
//			//更新配置节点信息
//			if event.Type == zk.EventNodeDataChanged {
//				content, _, chw, err := this.zkCon.GetW(path)
//				//content, _, err := this.zkCon.Get(path)
//				if err == nil {
//					configHandler(content)
//				}
//				ch = chw
//			}
//		}
//
//	})
//}
//
////func (this *ZKServiceCenter) AddServiceListener(listener service.Listener) {
////	this.Container.AddServiceListener(listener)
////}
