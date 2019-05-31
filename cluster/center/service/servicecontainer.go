/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/6/5
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package service

import (
	"sync"
)

func NewContainer() *Container {
	return &Container{root: make(map[string]*serviceCategory), lbs: make(map[string]string)}
}

type Container struct {
	sync.RWMutex
	root map[string]*serviceCategory //服务容器 key 服务名
	lbs  map[string]string           //服务的负载均衡策略
	//serviceListeners map[string]Listener  //服务监听
	//lbs string
}

func (this *Container) SetLbs(serviceName string, lbsStr string) {
	this.Lock()
	defer this.Unlock()
	this.lbs[serviceName] = lbsStr
	category := this.root[serviceName]
	if category != nil {
		category.setLbs(lbsStr)
	}
}

//func (this *Container) AddServiceListener(listener Listener) {
//	this.EnsureCategory(listener.GetServiceType()).AddListener(listener)
//}

//更新服务
func (this *Container) UpdateService(service IService, overwrite bool) bool {
	this.Lock()
	defer this.Unlock()
	category := this.EnsureCategory(service.GetName())
	result := category.updateService(service, overwrite)
	return result
}

func (this *Container) EnsureCategory(serviceName string) *serviceCategory {
	category := this.root[serviceName]

	if category == nil {
		category = NewServiceCategory(serviceName, this.lbs[serviceName])
		this.root[serviceName] = category
	}
	return category
}

func (this *Container) UpdateServices(serviceName string, services []IService) {
	this.Lock()
	defer this.Unlock()
	category := this.root[serviceName]
	if category == nil {
		category = NewServiceCategory(serviceName, this.lbs[serviceName])
		this.root[serviceName] = category
	}
	category.updateServices(services)

	//
	////TODO 关闭所有不可用的服务
	//for _, service := range services {
	//
	//
	//	//data, _, err := this.zkCon.Get(path + NODE_SPLIT + serviceID)
	//	//service := loadServiceFromData(data, serviceID, serviceName)
	//	//if service == nil {
	//	//	log.Errorf("%v unExpect service : %v", path, err)
	//	//	continue
	//	//}
	//	if category != nil {
	//		oldService := category.takeoutService(service)
	//		if oldService != nil {
	//			oldService.SetID(service.GetID())
	//			serviceCategory.updateService(oldService)
	//			continue
	//		}
	//	}
	//	//新服务需要连接上才能更新
	//	if service.Connect() {
	//		serviceCategory.updateService(service)
	//	}
	//}
	//this.root[serviceName] = serviceCategory
}

//删除服务
func (this *Container) RemoveService(serviceName string, serviceID string) {
	this.Lock()
	defer this.Unlock()
	serviceCategory := this.root[serviceName]
	if serviceCategory == nil {
		return
	}
	serviceCategory.removeService(serviceID)
}

//根据服务类型获取一个空闲的服务节点
func (this *Container) AllocService(serviceName string, param string) IService {
	this.Lock()
	defer this.Unlock()
	//TODO 后续要优化，考虑负载、空闲等因素
	serviceCategory := this.root[serviceName]
	if serviceCategory == nil {
		return nil
	}
	return serviceCategory.allocService(param)
}

//
//func (this *Container) GetMasterService(serviceType string) IService {
//	this.RLock()
//	defer this.RUnlock()
//	serviceCategory := this.root[serviceType]
//	if serviceCategory == nil {
//		return nil
//	}
//	return serviceCategory.getMaster()
//}

//获取指定服务节点
func (this *Container) GetService(serviceName string, serviceID string) IService {
	this.RLock()
	defer this.RUnlock()
	serviceCategory := this.root[serviceName]
	if serviceCategory == nil {
		return nil
	}
	return serviceCategory.services[serviceID]
}

func (this *Container) GetAllService(serviceType string) []IService {
	this.RLock()
	defer this.RUnlock()
	serviceCategory := this.root[serviceType]
	if serviceCategory == nil {
		return nil
	}
	return serviceCategory.getAllService()
}

func (this *Container) GetServiceInfo(serviceType string) []string {
	this.RLock()
	defer this.RUnlock()
	serviceCategory := this.root[serviceType]
	if serviceCategory == nil {
		return nil
	}
	return serviceCategory.getNodes()
}
