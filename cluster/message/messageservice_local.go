/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/5/6
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package message

////新建本地服务管理类
//func NewLocalService(serviceType string) *LocalService {
//	service := &LocalService{
//		open:        true,
//		dealTotal:   0,
//		serviceType: serviceType,
//		handlers:    make(map[int32]IServiceHandler),
//		//writeBack:   true,
//	}
//	return service
//}
//
////override IMessageService
//type LocalService struct {
//	sync.Mutex
//	open        bool  //服务是否开放
//	dealTotal   int64 //当前处理中的消息数量
//	serviceType string
//	handlers    map[int32]IServiceHandler
//	//writeBack   bool //是否回写消息管道
//}
//
//func (this *LocalService) GetType() string {
//	return this.serviceType
//}
//
////服务是否将相应消息回写
////func (this *LocalService) SetWriteBack(writeBack bool) {
////	this.writeBack = writeBack
////}
//
////注册消息服务处理句柄
//func (this *LocalService) RegisterHandler(seq int32, service IServiceHandler) {
//	this.handlers[seq] = service
//}
//
////是否能够处理指定编号的消息
//func (this *LocalService) CanHandle(seq int32) bool {
//	return this.handlers[seq] != nil
//}
//
////阻塞式调用消息接口
//func (this *LocalService) HandleMessage(request interface{}, out interface{}) error {
//	seqMessage, ok := request.(ISeqMessage)
//	if !ok {
//		return nil
//	}
//	if handler, ok := this.handlers[seqMessage.GetID()]; ok {
//		return handler.Request(request, out)
//	} else {
//		return nil
//	}
//}
//
////当前处理中的消息数据
//func (this *LocalService) GetDealTotal() int64 {
//	return this.dealTotal
//}
//
//func (this *LocalService) Close() {
//	this.open = false
//	timeout := 100 * time.Millisecond
//	//10次定时后还没处理完毕直接超时返回
//	for i := 0; i < 10; i++ {
//		time.Sleep(time.Duration(i) * timeout)
//		if this.dealTotal <= 0 {
//			return
//		}
//		log.Debugf("%v [%v] undeal message : %v", time.Now(), this.serviceType, this.dealTotal)
//	}
//}
