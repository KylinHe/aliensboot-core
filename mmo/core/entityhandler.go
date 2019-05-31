/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/12/8
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package core

import (
	"errors"
	"github.com/KylinHe/aliensboot-core/common/util"
)

type IEntityHandler interface {

	//持久化entity
	Save(entityID EntityID, entityType EntityType, data map[string]interface{}, callback func()) error

	//加载entity
	Load(entityID EntityID) (map[string]interface{}, error)

	//远程调用entity方法
	CallRemote(entityID EntityID, method string, args [][]byte) error

	//entity 迁移到远程节点
	MigrateRemote(spaceID EntityID, entityID EntityID, data []byte) error

	//获取定时器管理对象
	GetTimerManager() *util.TimerManager
}

var emptyHandler = &_EmptyHandler{}

type _EmptyHandler struct {
}

func (*_EmptyHandler) Save(entityID EntityID, entityType EntityType, data map[string]interface{}, callback func()) error {
	return errors.New("not implements")
}

func (*_EmptyHandler) Load(entityID EntityID) (map[string]interface{}, error) {
	return nil, errors.New("not implements")
}

func (*_EmptyHandler) CallRemote(entityID EntityID, method string, args [][]byte) error {
	return errors.New("not implements")
}

func (*_EmptyHandler) MigrateRemote(spaceID EntityID, entityID EntityID, data []byte) error {
	return errors.New("not implements")
}

func (*_EmptyHandler) GetTimerManager() *util.TimerManager {
	return nil
}
