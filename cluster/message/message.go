package message

import "errors"

var invalidServiceError error = errors.New("invalid service")

//消息服务抽象层，可以为local也可以是remote
type IMessageService interface {
	GetType() string                                        //获取消息服务类型
	HandleMessage(request interface{}) (interface{}, error) //阻塞调用消息服务接口
}

//type IRemoteService interface {
//	GetType() string                          //获取消息服务类型
//	BroadcastAll(message interface{}) //广播所有消息服务
//	RequestNode(serviceID string, request interface{}) (interface{}, error)      //阻塞调用指定服务ID接口，
//	RequestPriorityNode(serviceID string, request interface{}) (interface{}, error) //阻塞调用指定服务接口,优先发送到serviceID节点，没有会分配一个节点处理
//}

type ISeqMessage interface {
	GetID() int32 //获取消息id
}

type IServiceHandler interface {
	Request(request interface{}, out interface{}) error
}
