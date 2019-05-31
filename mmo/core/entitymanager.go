/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/8/31
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package core

import (
	"errors"
	"fmt"
	"github.com/KylinHe/aliensboot-core/common/data_structures/set"
	"github.com/KylinHe/aliensboot-core/common/util"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/KylinHe/aliensboot-core/mmo/unit"
	"github.com/vmihailenco/msgpack"
	"reflect"
	"time"
)

var EntityManager = newEntityManager()

type _EntityManager struct {
	handler IEntityHandler

	entities EntityMap //所有实体 id-entity

	entitiesByType map[EntityType]EntityMap //实体按类型分类 type[id-entity]

	entitiesDesc map[EntityType]*EntityDesc //实体元数据 type-entity_meta
}

func newEntityManager() *_EntityManager {
	return &_EntityManager{
		handler:        emptyHandler,
		entities:       EntityMap{},
		entitiesByType: map[EntityType]EntityMap{},
		entitiesDesc:   map[EntityType]*EntityDesc{},
	}
}

func (em *_EntityManager) addTimer(duration time.Duration, callbackFunc util.CallbackFunc) *util.Timer {
	return em.handler.GetTimerManager().AddTimer(duration, callbackFunc)
}

func (em *_EntityManager) addCallback(duration time.Duration, callbackFunc util.CallbackFunc) *util.Timer {
	return em.handler.GetTimerManager().AddCallback(duration, callbackFunc)
}

func (em *_EntityManager) put(entity *Entity) {
	em.entities.Add(entity)
	etype := entity.GetType()
	eid := entity.GetID()
	if entities, ok := em.entitiesByType[etype]; ok {
		entities.Add(entity)
	} else {
		em.entitiesByType[etype] = EntityMap{eid: entity}
	}
}

func (em *_EntityManager) del(e *Entity) {
	eid := e.GetID()
	em.entities.Del(eid)
	if entities, ok := em.entitiesByType[e.GetType()]; ok {
		entities.Del(eid)
	}
}

func (em *_EntityManager) traverseByType(eType EntityType, cb func(e *Entity)) {
	entities := em.entitiesByType[eType]
	for _, e := range entities {
		cb(e)
	}
}

// GenEntityID generates a new EntityID
func (em *_EntityManager) genEntityID() EntityID {
	return EntityID(util.GenUUID())
}

func (em *_EntityManager) RegisterHandler(handler IEntityHandler) {
	em.handler = handler
}

func (em *_EntityManager) GetHandler() IEntityHandler {
	return em.handler
}

func (em *_EntityManager) GetEntity(id EntityID) *Entity {
	return em.entities.Get(id)
}

//处理远程调用
func (em *_EntityManager) RemoteEntityCall(caller EntityID, id EntityID, method string, args [][]byte) (*Entity, error) {
	entity := em.GetEntity(id)
	if entity == nil {
		return nil, nil
	}
	return entity, entity.onCallFromRemote(caller, method, args)
}

//处理本地调用
func (em *_EntityManager) LocalEntityCall(id EntityID, method string, args []interface{}) (*Entity, error) {
	entity := em.GetEntity(id)
	if entity == nil {
		return nil, nil
	}
	return entity, entity.OnCallFromLocal(method, args)
}

//func (em *_EntityManager) UNRegisterEntity(entity IEntity) {
//	delete(em.entitiesDesc, typeName)
//}

// RegisterEntity registers custom entity type and define entity behaviors
func (em *_EntityManager) RegisterEntity(entity IEntity) *EntityDesc {
	entityVal := reflect.ValueOf(entity)
	entityType := entityVal.Type()
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}
	typeName := EntityType(entityType.Name())

	if desc, ok := em.entitiesDesc[typeName]; ok {
		log.Warnf("RegisterEntity: Entity type %s already registered", typeName)
		return desc
	}

	methodDesc := methodDescMap{}
	entityTypeDesc := &EntityDesc{
		name:            typeName,
		useAOI:          false,
		entityType:      entityType,
		methodDesc:      methodDesc,
		clientAttrs:     set.StringSet{},
		allAttrs:        set.StringSet{},
		persistAttrs:    set.StringSet{},
		persistInterval: time.Minute, //
	}
	em.entitiesDesc[typeName] = entityTypeDesc

	entityPtrType := reflect.PtrTo(entityType)
	numMethods := entityPtrType.NumMethod()

	for i := 0; i < numMethods; i++ {
		method := entityPtrType.Method(i)
		methodDesc.visit(method)
	}

	//// define entity Attrs
	entity.DescribeEntityType(entityTypeDesc)
	log.Infof(">>> RegisterEntity %s => %s <<<", typeName, entityType.Name())
	return entityTypeDesc
}

//从元数据中初始化一个实体
func (em *_EntityManager) CreateEntity(entityType EntityType, space *Space, pos unit.Vector, entityID EntityID) (*Entity, error) {
	entityDesc, ok := em.entitiesDesc[entityType]

	if !ok {
		return nil, errors.New(fmt.Sprintf("unknown entity type: %s", entityType))
	}

	//没有实体id自定义生成
	if entityID == "" {
		entityID = em.genEntityID()
	}

	entityInstance := reflect.New(entityDesc.entityType)
	entity := reflect.Indirect(entityInstance).FieldByName("Entity").Addr().Interface().(*Entity)
	//entity := &entity1
	entity.desc = entityDesc

	entity.init(entityID, entityInstance)
	em.put(entity)

	log.Debugf("Entity %s created.", entity)

	entity.I.OnCreated()

	if space != nil {
		space.enter(entity, pos)
	}

	// startup the periodical timer for saving entity
	if entity.IsPersistent() {
		entity.setupSaveTimer()
	}

	return entity, nil
}

//
func (em *_EntityManager) MigrateOut(spaceID EntityID, entityID EntityID) error {
	entity := em.GetEntity(entityID)
	if entity == nil {
		return errors.New(fmt.Sprintf("migrate entity not found : %v", entityID))
	}
	if entity.space.GetID() == spaceID {
		return errors.New(fmt.Sprintf("migrate entity already exist : %v", entityID))
	}

	migrateData := entity.GetMigrateData()
	data, err := msgpack.Marshal(migrateData)
	if err != nil {
		return errors.New(fmt.Sprintf("%s is migrating to space %s, but pack migrate data failed: %s", entity, spaceID, err))
	}
	// disable the entity
	entity.destroyEntity(true)
	em.handler.MigrateRemote(spaceID, entityID, data)
	return nil
}

//
func (em *_EntityManager) MigrateIn(entityID EntityID, space *Space, data []byte) error {
	var mData *entityMigrateData
	err := msgpack.Unmarshal(data, &mData)
	if err != nil {
		return err
	}

	typeName := mData.Type
	entityTypeDesc, ok := em.entitiesDesc[typeName]
	if !ok {
		return errors.New(fmt.Sprintf("restore unknown entity type: %s", typeName))
	}

	var entity *Entity
	var entityInstance reflect.Value

	entityInstance = reflect.New(entityTypeDesc.entityType)

	entity = reflect.Indirect(entityInstance).FieldByName("Entity").Addr().Interface().(*Entity)
	entity.desc = entityTypeDesc
	entity.init(entityID, entityInstance)

	entity.Position = mData.Pos
	entity.Yaw = mData.Yaw

	em.put(entity)
	entity.AssignMap(mData.Attrs)

	if space != nil {
		space.enter(entity, mData.Pos)
	}

	entity.I.OnAttrsReady()
	entity.I.OnMigrateIn()

	isPersistent := entity.desc.IsPersistent()
	if isPersistent { // startup the periodical timer for saving e
		entity.setupSaveTimer()
	}

	timerData := mData.TimerData
	if timerData != nil {
		entity.restoreTimers(timerData)
	}

	return nil
}
