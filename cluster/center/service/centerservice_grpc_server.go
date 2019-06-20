/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/8/20
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package service

import (
	"github.com/KylinHe/aliensboot-core/chanrpc"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/KylinHe/aliensboot-core/protocol/base"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"time"
)

const (
	suspendedTimeOut = time.Millisecond * 500
	commandRequest   = "request"
	commandReceive   = "receive"
)

type handler func(request *base.Any) *base.Any

func NewRpcHandler(chanRpc *chanrpc.Server, handler handler) *rpcServer {
	if chanRpc == nil {
		log.Fatalf("chanRpc can not be nil")
	}
	service := &rpcServer{}
	service.chanRpc = chanRpc
	service.handler = handler
	service.chanRpc.Register(commandRequest, service.request)
	service.chanRpc.Register(commandReceive, service.receive)
	return service
}

type rpcServer struct {
	chanRpc   *chanrpc.Server
	handler   handler
	suspended bool
	//启动服务参数
	server *grpc.Server //
}

func (this *rpcServer) start(name string, port int) bool {
	server := grpc.NewServer()
	base.RegisterRPCServiceServer(server, this)
	address := ":" + strconv.Itoa(port)
	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Errorf("failed to listen: %v", err)
		return false
	}
	go func() {
		_ = server.Serve(lis)
		//log.Infof("rpc service %v stop", name)
	}()
	this.server = server
	return true
}

func (this *rpcServer) close() {
	if this.server != nil {
		this.server.Stop()
	}
}

func (this *rpcServer) request(args []interface{}) {
	request := args[0].(*base.Any)
	server := args[1].(base.RPCService_RequestServer)
	response := this.handler(request)
	if response != nil {
		response.Id = request.Id
		_ = server.Send(response)
	}
}

func (this *rpcServer) receive(args []interface{}) {
	request := args[0].(*base.Any)
	this.handler(request)
}

func (this *rpcServer) LocalRequest(request *base.Any) (*base.Any, error) {
	return this.handler(request), nil
}


func (this *rpcServer) Request(request *base.Any, server base.RPCService_RequestServer) error {
	return this.chanRpc.Call0(commandRequest, request, server)
}

func (this *rpcServer) Receive(server base.RPCService_ReceiveServer) error {
	for {
		if this.suspended {
			time.Sleep(suspendedTimeOut)
			continue
		}
		request, err := server.Recv()
		if err != nil {
			//log.Debugf("accept async message error : %v", err)
			return err
		}
		this.chanRpc.Go(commandReceive, request)
	}
	return nil
}
