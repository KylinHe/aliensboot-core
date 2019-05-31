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
	"context"
	"github.com/KylinHe/aliensboot-core/common/util"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/KylinHe/aliensboot-core/protocol/base"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type HttpService struct {
	*CenterService
	srv *http.Server     `json:"-"`
	handler *gin.Engine  `json:"-"`
}

func (this *HttpService) GetDesc() string {
	return "http service"
}

//启动服务
func (this *HttpService) Start() bool {
	this.srv = &http.Server{
		Addr:   ":" + util.IntToString(this.Port),
		Handler: this.handler,
	}
	go func() {
		log.Debugf("Http Bind Port %v", this.srv.Addr)
		if err := this.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("start http service err : %v", err)
			this.srv = nil
		}
	}()
	return true
}

//连接服务
func (this *HttpService) Connect() bool {
	return true
}

//关闭服务
func (this *HttpService) Close() {
	//log.Debug("Shutdown Server ...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := this.srv.Shutdown(ctx); err != nil {
		log.Fatal("Http Server Shutdown:", err)
	}
	//log.Debug("Server exit")
	this.srv = nil
}

func (this *HttpService) SetHandler(handler interface{}) {
	result, ok := handler.(*gin.Engine)
	if !ok {
		log.Fatalf("invalid http service handler%v", reflect.TypeOf(handler))
	}
	this.handler = result
}

//比较服务是否冲突
func (this *HttpService) Equals(other IService) bool {
	otherService, ok := other.(*HttpService)
	if !ok {
		return false
	}
	return this.Name == otherService.Name && this.Address == otherService.Address && this.Port == otherService.Port
}

//服务是否本进程启动的
func (this *HttpService) IsLocal() bool {
	return this.srv != nil
}

//向服务请求消息
func (this *HttpService) Request(request *base.Any) (*base.Any, error) {
	return nil, nil
}

func (this *HttpService) Send(request *base.Any) error {
	return nil
}




