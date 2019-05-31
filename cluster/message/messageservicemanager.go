///*******************************************************************************
// * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
// * All rights reserved.
// * Date:
// *     2017/7/27
// * Contributors:
// *     aliens idea(xiamen) Corporation - initial API and implementation
// *     jialin.he <kylinh@gmail.com>
// *******************************************************************************/
package message

//
//import (
//	"sync"
//	"github.com/KylinHe/aliensboot-core/log"
//	"errors"
//)
//var invalidServiceError error = errors.NewTimeWheel("invalid service")
//
//var ServiceManager = initMessageServiceManager()
//
////var INVALID_SERVICE_RESPONSE = &protocol.GS2C{
////	Sequence: []int32{EXCEPTION_PUSH_SEQ},
////	ResultPush: &protocol.ResultPush{
////		Result: proto.Int32(int32(exception.SERVICE_INVALID)),
////	},
////}
//
//func initMessageServiceManager() *MessageServiceManager {
//	manager := &MessageServiceManager{
//		services: make(map[string]IMessageService),
//	}
//	return manager
//}
//
////服务容器,管理本地加载的服务句柄
//type MessageServiceManager struct {
//	sync.RWMutex
//	services map[string]IMessageService //模块业务服务句柄  处理业务消息
//}
//
////注册模块消息服务
//func (this *MessageServiceManager) RegisterService(service IMessageService) {
//	this.Lock()
//	defer this.Unlock()
//	log.Debug("register service %v", service.GetType())
//	this.services[service.GetType()] = service
//}
//
//func (this *MessageServiceManager) UnRegisterService(service IMessageService) {
//	this.Lock()
//	defer this.Unlock()
//	log.Debug("unregister service %v", service.GetType())
//	delete(this.services, service.GetType())
//}
//
//
////处理模块消息
//func (this *MessageServiceManager) Request(serviceType string, request interface{}) error, out interface{} {
//	service := this.services[serviceType]
//	if service != nil {
//		return service.Request(in, out)
//	}
//	return invalidServiceError
//}
//
////处理远程模块消息
//func (this *MessageServiceManager) RequestNode(serviceType string, serviceID string, request interface{}, out interface{}) error {
//	service := this.services[serviceType]
//	if service != nil {
//		remoteService, ok := service.(IRemoteService)
//		if ok {
//			return remoteService.RequestNode(serviceID, in, out)
//		}
//	}
//	return invalidServiceError
//}
//
////优先发送到指定的serviceID,如果没有发送到其他节点
//func (this *MessageServiceManager) RequestPriorityNode(serviceType string, serviceID string, request interface{}, out interface{}) error {
//	service := this.services[serviceType]
//	if service != nil {
//		remoteService, ok := service.(IRemoteService)
//		if ok {
//			return remoteService.RequestPriorityNode(serviceID, in, out)
//		}
//	}
//	return invalidServiceError
//}
//
//
//func (this *MessageServiceManager) BroadcastAllRemote(serviceType string, message interface{}) bool {
//	service := this.services[serviceType]
//	if service != nil {
//		remoteService, ok := service.(IRemoteService)
//		if ok {
//			return remoteService.BroadcastAll(message)
//		}
//	}
//	return false
//}
