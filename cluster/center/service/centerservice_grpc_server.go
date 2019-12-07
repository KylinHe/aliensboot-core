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

type Processor interface {
	NewResponseData() (interface{}, error)
	Decode(buf []byte) (interface{}, error)
	Encode(data interface{}) ([]byte, error)
}

type handler func(ctx *Context)

func NewRpcHandler(chanRpc *chanrpc.Server, handler handler, msgProcessor Processor) *rpcServer {
	service := &rpcServer{}
	service.chanRpc = chanRpc
	service.handler = handler
	service.msgProcessor = msgProcessor

	if chanRpc != nil {
		chanRpc.Register(commandRequest, service.request)
		chanRpc.Register(commandReceive, service.receive)
	}
	return service
}

type rpcServer struct {
	chanRpc   *chanrpc.Server
	handler   handler
	suspended bool
	//启动服务参数
	server *grpc.Server //
	msgProcessor Processor
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
	ctx := args[0].(*Context)
	this.handler(ctx)
}

func (this *rpcServer) receive(args []interface{}) {
	ctx := args[0].(*Context)
	this.handler(ctx)
}

//func (this *rpcServer) LocalRequest(request *base.Any) (*base.Any, error) {
//	return this.handler(request), nil
//}


func (this *rpcServer) Request(request *base.Any, server base.RPCService_RequestServer) error {
	ctx, err := this.newContext(request)
	if err != nil {
		return err
	}
	if this.chanRpc != nil {
		err := this.chanRpc.Call0(commandRequest, ctx, request)
		if err != nil {
			return err
		}
	} else {
		this.handler(ctx)
	}
	if ctx.ret {
		data , err := this.msgProcessor.Encode(ctx.Response)
		if err != nil {
			return err
		}
		ctx.response.Value = data
		server.Send(ctx.response)
	}
	return nil
}

func (this *rpcServer) Receive(server base.RPCService_ReceiveServer) error {
	for {
		if this.suspended {
			time.Sleep(suspendedTimeOut)
			continue
		}
		request, err := server.Recv()
		if err != nil {
			return err
		}
		ctx, err := this.newContext(request)
		if err != nil {
			return err
		}
		if this.chanRpc != nil {
			this.chanRpc.Go(commandReceive, ctx)
		} else {
			this.handler(ctx)
		}
	}
	return nil
}

func (this *rpcServer) newContext(request *base.Any) (*Context, error) {
	requestProxy, err := this.msgProcessor.Decode(request.Value)
	if err != nil {
		return nil, err
	}
	responseProxy, err := this.msgProcessor.NewResponseData()
	if err != nil {
		return nil, err
	}
	ctx := newContext(request, requestProxy, responseProxy)
	return ctx, nil
}
