package service

import (
	"errors"
	"github.com/KylinHe/aliensboot-core/protocol/base"
	"github.com/gogo/protobuf/proto"
)

func newContext(request *base.Any, server base.RPCService_RequestServer) *Context {
	return &Context{
		Request:request,
		Response:&base.Any{Id:request.Id},
		server:server,
		autoResp:true,
	}
}

type Context struct {

	Request *base.Any  // 请求消息

	Response *base.Any  // 响应消息

	server base.RPCService_RequestServer // 写消息句柄

	autoResp bool // 自动响应请求
}

// 上下验权通过
func (ctx *Context) Auth(authID int64) {
	ctx.Response.AuthId = authID
}

func (ctx *Context) SetAutoResp(auto bool) {
	ctx.autoResp = auto
}

func (ctx *Context) IsAutoResp() bool {
	return ctx.autoResp
}

// 响应proto消息
func (ctx *Context) GOGOProto(msg proto.Message) error {
	if msg == nil {
		return errors.New("msg can not be nil")
	}
	if ctx.server == nil {
		return errors.New("server not initial")
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	//any := &base.Any{}
	ctx.Response.Value = data
	// ctx.Response.AuthId = ctx.authID
	return ctx.server.Send(ctx.Response)
}

