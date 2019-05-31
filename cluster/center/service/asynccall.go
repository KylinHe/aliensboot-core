/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/10/16
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package service

import (
	"github.com/KylinHe/aliensboot-core/protocol/base"
)

func NewAsyncCall(request *base.Any, callProxy func(f func(), c func()), callbackProxy func(response *base.Any, err error)) *AsyncCall {
	return &AsyncCall{
		request:       request,
		callProxy:     callProxy,
		callbackProxy: callbackProxy,
	}
}

type AsyncCall struct {
	service IService

	request  *base.Any
	response *base.Any
	err      error

	callProxy     func(f func(), c func())
	callbackProxy func(response *base.Any, err error)
}

func (this *AsyncCall) ReqID() uint16 {
	return this.request.Id
}

//调用
func (this *AsyncCall) Invoke(service IService) {
	this.service = service
	this.callProxy(this.exec, this.callback)
}

func (this *AsyncCall) exec() {
	this.response, this.err = this.service.Request(this.request)
}

func (this *AsyncCall) callback() {
	if this.callbackProxy != nil {
		this.callbackProxy(this.response, this.err)
	}
}
