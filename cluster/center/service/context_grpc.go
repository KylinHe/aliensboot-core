package service

import (
	"github.com/KylinHe/aliensboot-core/protocol/base"
)

func newContext(request *base.Any, requestProxy interface{}, responseProxy interface{}) *Context {
	return &Context{
		Request: requestProxy,
		Response: responseProxy,
		request:request,
		response:&base.Any{Id:request.Id},
	}
}

type Context struct {

	Request interface{}

	Response interface{}

	request *base.Any  // 请求消息

	response *base.Any  // 响应消息

	ret bool
}

// 获取消息id
func (ctx *Context) GetMsgId() uint16 {
	return ctx.request.Id
}

// 获取序号id
func (ctx *Context) GetSeqId() uint32 {
	return ctx.request.SeqId
}

// 获取权限id
func (ctx *Context) GetAuthId() int64 {
	return ctx.request.AuthId
}

// 获取网关id
func (ctx *Context) GetGateID() string {
	return ctx.request.GateId
}

func (ctx *Context) GetHeader(key string) []byte {
	return ctx.request.GetHeaderByKey(key)
}

func (ctx *Context) GetHeaderStr(key string) string {
	return ctx.request.GetHeaderStrByKey(key)
}

// 上下验权通过
func (ctx *Context) Auth(authID int64) {
	if ctx.response != nil {
		ctx.response.AuthId = authID
	}
}

//func (ctx *Context) SetAutoResp(auto bool) {
//	ctx.autoResp = auto
//}
//
//func (ctx *Context) IsAutoResp() bool {
//	return ctx.autoResp
//}

// 响应消息
func (ctx *Context) WriteResponse() {
	ctx.ret = true
	//if ctx.server == nil {
	//	return errors.New("server not initial")
	//}
	//data , err := ctx.msgProcessor.Encode(ctx.Response)
	//if err != nil {
	//	return err
	//}
	//ctx.response.Value = data
	//return ctx.server.Send(ctx.response)
}

