/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2019/12/2
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package center

import "github.com/KylinHe/aliensboot-core/module/base"

func NewPrefixDataProxy(id string, skeleton *base.Skeleton, listener DataPrefixListener) *_prefixDataProxy {
	proxy := &_prefixDataProxy{
		id: id,
		skeleton: skeleton,
		listener: listener,
	}
	skeleton.RegisterChanRPC(id, proxy.handle)
	return proxy
}

type _prefixDataProxy struct {
	id string
	skeleton *base.Skeleton
	listener DataPrefixListener
}

func (proxy *_prefixDataProxy) OnDataChange(eventType DataEventType, data []byte, dataRootName string, dataName string, init bool) {
	if init {
		proxy.listener(eventType, data, dataRootName, dataName, init)
	} else {
		_ = proxy.skeleton.ChanRPCServer.Call0(proxy.id, eventType, data, dataRootName, dataName, init)
	}
	//_ = proxy.skeleton.ChanRPCServer.Call0(proxy.id, eventType, data, dataRootName, dataName)
}

func (proxy *_prefixDataProxy) handle(param []interface{}) {
	eventType := param[0].(DataEventType)
	data := param[1].([]byte)
	dataRootName := param[2].(string)
	dataName := param[3].(string)
	init := param[4].(bool)
	proxy.listener(eventType, data, dataRootName, dataName, init)
}


func NewDataProxy(id string, skeleton *base.Skeleton, listener DataListener) *_dataProxy {
	proxy := &_dataProxy{
		id: id,
		skeleton: skeleton,
		listener: listener,
	}
	skeleton.RegisterChanRPC(id, proxy.handle)
	return proxy
}

type _dataProxy struct {
	id string
	skeleton *base.Skeleton
	listener DataListener
}

func (proxy *_dataProxy) OnDataChange(content []byte, init bool) {
	if init {
		proxy.listener(content, init)
	} else {
		_ = proxy.skeleton.ChanRPCServer.Call0(proxy.id, content, init)
	}
}

func (proxy *_dataProxy) handle(param []interface{}) {
	content := param[0].([]byte)
	init := param[1].(bool)
	proxy.listener(content, init)
}
