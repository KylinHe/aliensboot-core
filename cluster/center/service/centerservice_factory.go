/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/6/4
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package service

import (
	"errors"
	"github.com/KylinHe/aliensboot-core/cluster/center/lbs"
	"github.com/KylinHe/aliensboot-core/config"
)


func NewService(config config.ServiceConfig) (IService, error) {
	if !lbs.ValidateLBS(config.Lbs) {
		return nil, errors.New("unexpect lbs option " + config.Lbs)
	}
	return NewService1(config.ID, config.Name, config.Address, config.Port, config.Protocol)
}

func NewService2(centerService *CenterService, id string, name string) (IService, error) {
	centerService.SetID(id)
	centerService.SetName(name)
	var service IService = nil
	switch centerService.Protocol {
		case GRPC:
			service = &GRPCService{CenterService: centerService}
			break
		//case WEBSOCKET:
		//	return &wbService{centerService: centerService}
		case HTTP:
			service = &HttpService{CenterService: centerService}
			break
	default:
		service = &GRPCService{CenterService: centerService}

	}

	return service, nil
}

func NewService1(id string, name string, address string, port int, protocol string) (IService, error) {
	centerService := &CenterService{
		Address:  address,
		Port:     port,
		Protocol: protocol,
	}
	return NewService2(centerService, id, name)
}
