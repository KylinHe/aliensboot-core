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

//func PublicWBService(config CenterService, address string) *wbService {
//	if !ClusterCenter.IsConnect() {
//		panic(config.Name + " cluster center is not connected")
//		return nil
//	}
//	service := &wbService{
//		&centerService{
//			id:       ClusterCenter.GetNodeID(),
//			name:     config.Name,
//			Address:  address,
//			Protocol: WEBSOCKET,
//		},
//	}
//	//center.ClusterCenter.AddServiceFactory(service.name, &wbServiceFactory{})
//	//websocket服务启动成功,则发布到中心服务器
//	if !ClusterCenter.PublicService(service, config.Unique) {
//		panic(service.name + " wb service can not be public")
//	}
//	return service
//}
//
//type wbService struct {
//	*centerService
//}
//
//func (this *wbService) GetDesc() string {
//	return "websocket service"
//}
//
//func (this *wbService) GetID() string {
//	return this.id
//}
//
//func (this *wbService) GetType() string {
//	return this.name
//}
//
//func (this *wbService) SetID(id string) {
//	this.id = id
//}
//
//func (this *wbService) SetType(serviceType string) {
//	this.name = serviceType
//}
//
////启动服务
//func (this *wbService) Start() bool {
//	return true
//}
//
////连接服务
//func (this *wbService) Connect() bool {
//	return true
//}
//
////比较服务是否冲突
//func (this *wbService) Equals(other IService) bool {
//	otherService, ok := other.(*wbService)
//	if !ok {
//		return false
//	}
//	return this.name == otherService.name && this.Address == otherService.Address
//}
//
////服务是否本进程启动的
//func (this *wbService) IsLocal() bool {
//	return true
//}
//
////关闭服务
//func (this *wbService) Close() {
//}
//
////向服务请求消息
//func (this *wbService) Request(request interface{}) (interface{}, error) {
//	return nil, nil
//}
