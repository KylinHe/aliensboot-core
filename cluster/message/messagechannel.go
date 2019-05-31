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

import (
	"github.com/KylinHe/aliensboot-core/exception"
	"github.com/KylinHe/aliensboot-core/log"
)

type MessageChannel struct {
	channel      chan interface{} //消息管道
	messageLimit int
	handler      func(msg interface{}) interface{}
}

//向管道发送消息
func (this *MessageChannel) WriteMsg(message interface{}) {
	//用户消息管道没开，不接受消息
	if this.channel == nil {
		return
	}
	select {
	case this.channel <- message:
	default:
		log.Debugf("message channel full %v - %v", this.channel, message)
		//TODO 消息管道满了需要异常处理
	}
}

func (this *MessageChannel) Open() {
	this.channel = make(chan interface{}, this.messageLimit)
	go func() {
		defer func() {
			exception.CatchStackDetail()
		}()
		for {
			//只要消息管道没有关闭，就一直等待用户请求
			message, open := <-this.channel
			if !open {
				break
			}
			this.handler(message)
		}
		this.Close()
	}()
}

//关闭消息管道
func (this *MessageChannel) Close() {
	if this.channel != nil {
		close(this.channel)
		this.channel = nil
	}
}
