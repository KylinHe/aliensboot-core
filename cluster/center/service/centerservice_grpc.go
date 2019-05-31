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

import (
	"github.com/KylinHe/aliensboot-core/common/util"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/KylinHe/aliensboot-core/protocol/base"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"reflect"
	"time"
)

const (
	requestTimeout = 10 * time.Second
)

type GRPCService struct {
	*CenterService

	//pid *actor.PID

	server *rpcServer //grpc服务端处理句柄

	client *grpc.ClientConn

	requestClient base.RPCServiceClient
	receiveClient base.RPCService_ReceiveClient
}

func (this *GRPCService) GetConfig() *CenterService {
	return this.CenterService
}

func (this *GRPCService) GetDesc() string {
	return "rpc service"
}

func (this *GRPCService) SetHandler(handler interface{}) {
	result, ok := handler.(*rpcServer)
	if !ok {
		log.Fatalf("invalid grpc service handler %v", reflect.TypeOf(handler))
	}
	this.server = result
}

//启动服务
func (this *GRPCService) Start() bool {
	//var server *grpc.Server = nil
	//if ClusterCenter.certFile != "" {
	//	creds, err := credentials.NewServerTLSFromFile(ClusterCenter.certFile, ClusterCenter.keyFile)
	//	if err != nil {
	//		log.Fatalf("Failed to generate credentials %v", err)
	//		return false
	//	}
	//	server = grpc.NewServer(grpc.Creds(creds))
	//} else {
	//	server = grpc.NewServer()
	//}
	return this.server.start(this.Name, this.Port)
}

//连接服务
func (this *GRPCService) Connect() bool {
	var option grpc.DialOption = nil
	//if ClusterCenter.certFile != "" {
	//	creds, err := credentials.NewClientTLSFromFile(ClusterCenter.certFile, ClusterCenter.commonName)
	//	if err != nil {
	//		log.Errorf("Failed to create TLS credentials %v", err)
	//		return false
	//	}
	//	option = grpc.WithTransportCredentials(creds)
	//} else {
	//	option = grpc.WithInsecure()
	//}
	option = grpc.WithInsecure()
	address := this.Address + ":" + util.IntToString(this.Port)
	conn, err := grpc.Dial(address, option)
	if err != nil {
		log.Errorf("did not connect: %v", err)
		return false
	}
	this.client = conn
	this.requestClient = base.NewRPCServiceClient(this.client)

	receiveClient, err := this.requestClient.Receive(context.Background())
	if err != nil {
		log.Errorf("init receive handler error : %v", err)
		return false
	}
	this.receiveClient = receiveClient

	log.Infof("connect grpc service %v-%v success", this.Name, address)
	return true
}

//比较服务是否冲突
func (this *GRPCService) Equals(other IService) bool {
	otherService, ok := other.(*GRPCService)
	if !ok {
		return false
	}
	return this.Name == otherService.Name && this.Address == otherService.Address && this.Port == this.Port
}

//服务是否本进程启动的
func (this *GRPCService) IsLocal() bool {
	return this.server != nil
}

//关闭服务
func (this *GRPCService) Close() {
	if this.server != nil {
		this.server.close()
		this.server = nil
	}
	if this.client != nil {
		this.client.Close()
		this.client = nil
	}
}

//向服务请求消息
func (this *GRPCService) Request(request *base.Any) (*base.Any, error) {
	//本地启动的服务直接调用，不需要经过grpc中转
	if this.IsLocal() {
		return this.server.LocalRequest(request)
	}
	if this.requestClient == nil {
		return nil, errors.New("service is not initial")
	}

	//加入超时机制，防止卡死
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	client, err := this.requestClient.Request(ctx, request)
	if err != nil {
		return nil, err
	}
	return client.Recv()
}

//func (this *GRPCService) AsyncRequest(asyncCall *AsyncCall) {
//	//asyncCall.requestClient = this.requestClient
//	//asyncCall.Invoke()
//
//	//callback( , )
//	//this.pid.Tell(&call{request, callback})
//}

//让服务接收消息，不需要响应
func (this *GRPCService) Send(request *base.Any) error {
	//服务为本机，直接处理
	if this.IsLocal() {
		this.server.LocalRequest(request)
		return nil
	}

	if this.receiveClient == nil {
		return errors.New("service is not initial")
	}
	return this.receiveClient.Send(request)
}
