/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/4/28
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package service

import (
	"github.com/KylinHe/aliensboot-core/cluster/center/lbs"
	"github.com/KylinHe/aliensboot-core/log"
	"reflect"
)

func NewServiceCategory(category string, lbsStr string) *serviceCategory {
	result := &serviceCategory{
		category: category,
		services: make(map[string]IService),
		lbsName: lbs.StrategyPolling,
		//nodes:    []string{},
		//listeners: []Listener{},
		//seqs:     seqMaps,
	}
	result.setLbs(lbsStr)
	return result
}

type serviceCategory struct {
	category string
	lbs      lbs.Strategy        //负载均衡策略
	lbsName  string
	services map[string]IService //服务节点名,和服务句柄
	//nodes    []string

	//listeners []Listener
	//seqs     map[int32]struct{} //能够处理的消息编号
}

func (this *serviceCategory) setLbs(lbsStr string) {
	newLbs := lbs.GetLBS(lbsStr)
	if this.lbs == nil || reflect.TypeOf(this.lbs) != reflect.TypeOf(newLbs) {
		this.lbs = newLbs
		this.lbsName = lbsStr
		log.Debugf("[lbs-%v] init %v", this.lbsName, this.category)
		for serviceId, _ := range this.services {
			this.handleAddNode(serviceId)
		}
	}
	//更新负载均衡策略

	//log.Debugf("update lbs strategy %v - %v", this.category, lbsStr)
}

//分配一个可用服务
func (this *serviceCategory) allocService(key string) IService {
	nodeName := this.lbs.GetNode(key)
	if nodeName == "" {
		return nil
	}
	return this.services[nodeName]
}

//func (this *serviceCategory) AddListener(listener Listener) {
//	this.listeners = append(this.listeners, listener)
//}

//func (this *serviceCategory) canHandle(messageSeq int32) bool {
//	_, ok := this.seqs[messageSeq]
//	return ok
//}

//初始化lbs节点信息
//func (this *serviceCategory) initLBSNode() {
//	nodes := []string{}
//	for node, _ := range this.services {
//		nodes = append(nodes, node)
//	}
//	sort.Strings(nodes)
//	this.nodes = nodes
//	this.lbs.Init(this.nodes)
//}

//更新服务
func (this *serviceCategory) updateService(service IService, overwrite bool) bool {
	//不允许覆盖存在直接返回
	oldService, _ := this.services[service.GetID()]
	if !overwrite {
		if oldService != nil {
			return false
		}
	}

	if oldService != nil {
		oldService.Close()
	}

	if !service.IsLocal() {
		service.Connect()
	}

	this.services[service.GetID()] = service
	this.handleAddNode(service.GetID())
	return true
}

//取出相同的服务
func (this *serviceCategory) updateServices(services []IService) {
	newServices := make(map[string]IService) //服务节点名,和服务句柄
	for _, service := range services {
		oldService := this.services[service.GetID()]
		if oldService != nil {
			newServices[service.GetID()] = oldService
			log.Debugf("repeated service %v", oldService)
			delete(this.services, service.GetID())
		} else if service.Connect() {
			newServices[service.GetID()] = service
			log.Debugf("new connect service %v", service)
			this.handleAddNode(service.GetID())
		}
	}
	for _, releaseService := range this.services {
		releaseService.Close()
		this.handleRemoveNode(releaseService.GetID())
	}
	this.services = newServices
	//服务地址信息没有变，不需要再连接
	//for key, service := range this.services {
	//	if service.Equals(serviceConfig) {
	//		delete(this.services, key)
	//		this.initLBSNode()
	//		return service
	//	}
	//}
	//服务地址信息没有变，不需要再连接
	//for key, service := range this.services {
	//	if service.Equals(serviceConfig) {
	//		delete(this.services, key)
	//		this.initLBSNode()
	//		return service
	//	}
	//}
	//return nil

}

func (this *serviceCategory) removeService(serviceID string) {
	removeService, ok := this.services[serviceID]
	if !ok {
		return
	}
	if removeService != nil {
		removeService.Close()
	}
	delete(this.services, serviceID)
	this.handleRemoveNode(serviceID)
}

func (this *serviceCategory) handleRemoveNode(serviceID string) {
	this.lbs.RemoveNode(serviceID)
	log.Debugf("[lbs-%v] remove node %v-%v", this.lbsName, this.category, serviceID)
}

func (this *serviceCategory) handleAddNode(serviceID string) {
	this.lbs.AddNode(serviceID, 1)
	log.Debugf("[lbs-%v] add node %v-%v", this.lbsName, this.category, serviceID)
}

//func (this *serviceCategory) getNodes() []string {
//	return this.nodes
//}

func (this *serviceCategory) getAllService() []IService {
	results := []IService{}
	for _, service := range this.services {
		results = append(results, service)
	}
	return results
}

//func (this *serviceCategory) getMaster() IService {
//	//TODO 后续要加一套master-salve机制
//	if len(this.nodes) == 0 {
//		return nil
//	}
//	return this.services[this.nodes[0]]
//}
