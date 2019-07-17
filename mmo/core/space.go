package core

import (
	"fmt"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/KylinHe/aliensboot-core/mmo/aoi"
	"github.com/KylinHe/aliensboot-core/mmo/config"
	"github.com/KylinHe/aliensboot-core/mmo/unit"
	"time"
)

type ISpace interface {
	IEntity

	// Called when initializing space struct, override to initialize custom space fields
	OnSpaceInit()

	// Called when space is created
	OnSpaceCreated()

	// Called just before space is destroyed
	OnSpaceDestroy()

	// Called when any entity enters space
	OnEntityEnterSpace(entity *Entity)

	// Called when any entity leaves space
	OnEntityLeaveSpace(entity *Entity)
}

const (
	SpaceAttrType = "spaceType"

	tickInterval = 100 * time.Millisecond
)

type Space struct {
	Entity

	I ISpace

	entities EntitySet //entities in current space

	aoiMgr aoi.Manager
}

// OnInit initialize Space entity
func (space *Space) OnInit() {
	space.entities = EntitySet{}
	space.I = space.Entity.I.(ISpace)

	spaceConfig := &config.SpaceConfig{"testSpace", -5000, 5000, -5000, 5000, 500}
	space.aoiMgr = aoi.NewTowerAOIManager(spaceConfig.MinX, spaceConfig.MaxX, spaceConfig.MinY, spaceConfig.MaxY, spaceConfig.TowerRange)
	space.I.OnSpaceInit()

	space.addRawTimer(tickInterval, space.OnTick)
}

//tick
func (space *Space) OnTick([]interface{}) {
	for entity, _ := range space.entities {
		entity.OnTick(tickInterval)
	}
}

func (space *Space) OnCreated() {
	space.I.OnSpaceCreated()
	SpaceManager.putSpace(space)
}

// OnDestroy is called when Space entity is destroyed
func (space *Space) OnDestroy() {
	space.I.OnSpaceDestroy()
	for e := range space.entities {
		e.Destroy()
	}
	SpaceManager.delSpace(space.GetID())
}

func (space *Space) DescribeEntityType(desc *EntityDesc) {
	desc.DefineAttr(SpaceAttrType, AttrAllClient)
}

func (space *Space) String() string {
	return fmt.Sprintf("space<%d>", space.GetID())
}

//进入场景
func (space *Space) enter(entity *Entity, pos unit.Vector) {
	if entity.space != nil {
		log.Panicf("%s.enter(%s): current space is not nil, but %s", space, entity, entity.space)
	}

	entity.space = space
	space.entities.Add(entity)

	if space.aoiMgr != nil && entity.IsUseAOI() {
		space.aoiMgr.Enter(entity.aoi, pos.X, pos.Y)
	}

	entity.I.OnEnterSpace()
	space.I.OnEntityEnterSpace(entity)
}

//离开场景
func (space *Space) leave(entity *Entity) {
	if entity.space != space {
		log.Panicf("%s.leave(%s): entity is not in this Space", space, entity)
	}

	entity.space = nil
	space.entities.Del(entity)
	if space.aoiMgr != nil && entity.IsUseAOI() {
		space.aoiMgr.Leave(entity.aoi)
	}
	space.I.OnEntityLeaveSpace(entity)
	entity.I.OnLeaveSpace(space)
}

//场景中移动
func (space *Space) move(entity *Entity, newPos unit.Vector) {
	if entity.space != space {
		log.Panicf("%s.leave(%s): entity is not in this Space", space, entity)
	}

	entity.Position = newPos
	space.aoiMgr.Moved(entity.aoi, newPos.X, newPos.Y)
	//for neighbor, _ := range entity.interestedIn {
	//	neighbor.proxy.OnEntityMove(entity)
	//}
	//space.proxy.OnEntityMove(entity)
}

// CreateEntity creates a new local entity in this space
func (space *Space) CreateEntity(eType EntityType, pos unit.Vector, id EntityID) (*Entity, error) {
	return EntityManager.CreateEntity(eType, space, pos, id)
}

// GetEntityCount returns the total count of entities in space
func (space *Space) GetEntityCount() int {
	return len(space.entities)
}

// ForEachEntity visits all entities in space and call function f with each entity
func (space *Space) ForEachEntity(f func(e *Entity)) {
	for entity, _ := range space.entities {
		f(entity)
	}
}

//--------------------abstract method----------------

func (space *Space) OnSpaceInit() {

}

func (space *Space) OnSpaceCreated() {

}

func (space *Space) OnSpaceDestroy() {

}

func (space *Space) OnEntityEnterSpace(entity *Entity) {

}

func (space *Space) OnEntityLeaveSpace(entity *Entity) {

}
