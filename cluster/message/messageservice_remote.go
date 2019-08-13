/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/3/24
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package message

import (
	"errors"
	"github.com/KylinHe/aliensboot-core/cluster/center"
	"github.com/KylinHe/aliensboot-core/cluster/center/service"
	"github.com/KylinHe/aliensboot-core/protocol/base"
)

func NewRemoteService(serviceType string) *RemoteService {
	service := &RemoteService{
		serviceType: serviceType,
	}
	service.Init()
	return service
}

//远程调度服务 override IMessageService
type RemoteService struct {
	serviceType string //服务类型
}

func (this *RemoteService) Init() {
	center.ClusterCenter.SubscribeServices(this.serviceType)
}

//同步调用服务, 请求节点采用内部的负载均衡策略分配
func (this *RemoteService) Request(request *base.Any, param string) (*base.Any, error) {
	service := center.ClusterCenter.AllocService(this.serviceType, param)
	if service == nil {
		return nil, errors.New("invalid service:" + this.serviceType)
	}
	return service.Request(request)
}

//同步调用服务，请求到指定节点
func (this *RemoteService) RequestNode(serviceID string, request *base.Any) (*base.Any, error) {
	service := center.ClusterCenter.GetService(this.serviceType, serviceID)
	if service == nil {
		return nil, errors.New("invalid service:" + this.serviceType)
	}
	return service.Request(request)
}

//同步调用服务, 请求节点采用内部的负载均衡策略分配
func (this *RemoteService) AsyncRequest(param string, asyncCall *service.AsyncCall) error {
	service := center.ClusterCenter.AllocService(this.serviceType, param)
	if service == nil {
		return errors.New("invalid service:" + this.serviceType)
	}
	asyncCall.Invoke(service)
	//service.AsyncRequest(request, callback)
	return nil
}

//同步调用服务，请求到指定节点
func (this *RemoteService) AsyncRequestNode(serviceID string, asyncCall *service.AsyncCall) error {
	service := center.ClusterCenter.GetService(this.serviceType, serviceID)
	if service == nil {
		return errors.New("invalid service:" + this.serviceType)
	}
	asyncCall.Invoke(service)
	//service.AsyncRequest(request, callback)
	return nil
}

//func (this *RemoteService) RequestPriorityNode(serviceID string, request *base.Any) (*base.Any, error) {
//	service := center.ClusterCenter.GetService(this.serviceType, serviceID)
//	if service == nil {
//		service = center.ClusterCenter.AllocService(this.serviceType, "")
//	}
//	if service == nil {
//		return nil, invalidServiceError
//	}
//	return service.Request(request)
//}

//异步调用服务, 请求节点采用内部的负载均衡策略分配
func (this *RemoteService) Send(request *base.Any, param string) error {
	service := center.ClusterCenter.AllocService(this.serviceType, "")
	if service == nil {
		return errors.New("invalid service:" + this.serviceType)
	}
	return service.Send(request)
}

//异步调用服务，请求到指定节点
func (this *RemoteService) SendNode(serviceID string, request *base.Any) error {
	service := center.ClusterCenter.GetService(this.serviceType, serviceID)
	if service == nil {
		return errors.New("invalid service:" + this.serviceType)
	}
	return service.Send(request)
}

//广播到所有节点
func (this *RemoteService) BroadcastAll(message *base.Any) {
	services := center.ClusterCenter.GetAllService(this.serviceType)
	if services == nil || len(services) == 0 {
		return
	}
	for _, service := range services {
		_ = service.Send(message)
	}
	return
}

//获取消息服务类型
func (this *RemoteService) GetType() string {
	return this.serviceType
}
