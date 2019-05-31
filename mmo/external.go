/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/12/5
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package mmo

import (
	"errors"
	"fmt"
	"github.com/KylinHe/aliensboot-core/mmo/core"
	"github.com/KylinHe/aliensboot-core/mmo/unit"
	"github.com/vmihailenco/msgpack"
)

type PlayerClient interface {

	//CallRemote(id EntityID, method string, args []string)

}

type Space = core.Space

type Entity = core.Entity

type EntityID = core.EntityID

type EntityTimerID = core.EntityTimerID

type EntityType = core.EntityType

func RegisterSpace(spacePtr core.ISpace) {
	core.EntityManager.RegisterEntity(spacePtr)
}

func RegisterEntity(entity core.IEntity) {
	core.EntityManager.RegisterEntity(entity)
}

func RegisterEntityHandler(handler core.IEntityHandler) {
	core.EntityManager.RegisterHandler(handler)
}

func CreateSpace(eType EntityType, id EntityID) (*Space, error) {
	e, err := core.EntityManager.CreateEntity(eType, nil, unit.EmptyVector, id)
	return e.AsSpace(), err
}

func CreateEntity(eType EntityType, space *Space, pos unit.Vector) (*Entity, error) {
	return core.EntityManager.CreateEntity(eType, space, pos, "")
}

// GetSpace gets the space by ID
func GetSpace(id EntityID) *Space {
	return core.SpaceManager.GetSpace(id)
}

//entity迁移
func MigrateTo(spaceID EntityID, entityID EntityID) error {
	return core.EntityManager.MigrateOut(spaceID, entityID)
}

func MigrateIn(spaceID EntityID, entityID EntityID, data []byte) error {
	space := GetSpace(spaceID)
	if space == nil {
		return errors.New(fmt.Sprintf("space %v not found ", spaceID))
	}
	return core.EntityManager.MigrateIn(entityID, space, data)
}

//实体登录到场景
func EnterSpace(spaceID EntityID, eType EntityType, entityID EntityID, pos unit.Vector) (*Entity, error) {
	space := GetSpace(spaceID)
	if space == nil {
		return nil, errors.New(fmt.Sprintf("space %v not found ", spaceID))
	}
	entity := core.EntityManager.GetEntity(entityID)
	//实体已经存在
	if entity != nil {
		return entity, nil
	}
	return space.CreateEntity(eType, pos, entityID)
}

//handle
func RemoteEntityCall(caller EntityID, id EntityID, method string, args [][]byte) (*Entity, error) {
	return core.EntityManager.RemoteEntityCall(caller, id, method, args)
}

//call entity method
func Call(id EntityID, method string, args ...interface{}) error {
	entity, err := core.EntityManager.LocalEntityCall(id, method, args)
	//本地不存在、调用远程对象
	if entity == nil {
		argsData := make([][]byte, len(args))
		for i, arg := range args {
			data, err := msgpack.Marshal(arg)
			if err != nil {
				return err
			}
			argsData[i] = data
		}
		return core.EntityManager.GetHandler().CallRemote(id, method, argsData)
	}
	return err
}
