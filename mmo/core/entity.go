package core

import (
	"errors"
	"fmt"
	"github.com/KylinHe/aliensboot-core/common/util"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/KylinHe/aliensboot-core/mmo/aoi"
	"github.com/KylinHe/aliensboot-core/mmo/unit"
	"github.com/vmihailenco/msgpack"
	"reflect"
	"time"
	"unsafe"
)

type AI interface {
	OnTick(dt time.Duration)
}

// IEntity declares functions that is defined in Entity
// These functions are mostly component functions
type IEntity interface {

	// Called when entity attributes are ready.
	OnAttrsReady()

	//--------------Entity Lifetime----------

	OnInit()    // Called when initializing entity struct, override to initialize entity custom fields
	OnCreated() // Called when entity is just created
	OnDestroy() // Called when entity is destroying (just before destroy)

	// Migration
	OnMigrateOut() // Called just before entity is migrating out
	OnMigrateIn()  // Called justafter entity is migrating in

	// Space Operations
	OnEnterSpace()             // Called when entity leaves space
	OnLeaveSpace(space *Space) // Called when entity enters space

	DescribeEntityType(desc *EntityDesc) // Define entity attributes in this function

	String() string
}

type EntityType string

// EntityID type
type EntityID string

// IsNil returns if EntityID is nil
func (id EntityID) IsNil() bool {
	return id == ""
}

// entity info that should be migrated
type entityMigrateData struct {
	Type      EntityType             `msgpack:"T"`
	Attrs     map[string]interface{} `msgpack:"A"`
	Pos       unit.Vector            `msgpack:"Pos"`
	Yaw       unit.Yaw               `msgpack:"Yaw"`
	TimerData []byte                 `msgpack:"TD,omitempty"`
}

type Entity struct {
	*MapAttr //属性'

	desc *EntityDesc

	ai AI

	id EntityID

	//clientID ClientID

	I IEntity //实现类

	V reflect.Value

	Position unit.Vector

	Yaw unit.Yaw

	space *Space //实体所在的空间

	aoi *aoi.AOI //aoi

	interestedIn EntitySet //当前实体视野范围内的实体

	interestedBy EntitySet //视野内存在当前实体的对象

	destroyed bool

	//handle timer
	lastTimerId EntityTimerID
	timers      map[EntityTimerID]*entityTimerInfo
	rawTimers   map[*util.Timer]struct{}
}

func (e *Entity) GetID() EntityID {
	return e.id
}

func (e *Entity) GetSpaceID() EntityID {
	if e.space == nil {
		return ""
	}
	return e.space.id
}

func (e *Entity) GetType() EntityType {
	return e.desc.name
}

// SetPosition sets the entity position
func (e *Entity) SetPosition(pos unit.Vector) {
	e.setPositionYaw(pos, e.Yaw)
}

func (e *Entity) IsUseAOI() bool {
	return e.desc.useAOI
}

func (e *Entity) setPositionYaw(pos unit.Vector, yaw unit.Yaw) {
	space := e.space
	if space == nil {
		log.Warnf("%s.SetPosition(%s): space is nil", e, pos)
		return
	}
	space.move(e, pos)
	e.Yaw = yaw
}

func (e *Entity) init(entityID EntityID, entityInstance reflect.Value) {
	e.id = entityID
	e.V = entityInstance
	e.I = entityInstance.Interface().(IEntity)

	attrs := NewMapAttr()
	attrs.owner = e
	e.MapAttr = attrs
	e.timers = make(map[EntityTimerID]*entityTimerInfo)
	e.rawTimers = make(map[*util.Timer]struct{})

	e.interestedIn = make(EntitySet)
	e.interestedBy = make(EntitySet)
	e.aoi = aoi.NewAOI(e, e.desc.aoiDistance)

	e.I.OnInit()
}

func (e *Entity) OnTick(duration time.Duration) {
	if e.ai != nil {
		e.ai.OnTick(duration)
	}
}

func (e *Entity) SetAI(ai AI) {
	e.ai = ai
}

func (e *Entity) OnEnterAOI(otherAoi *aoi.AOI) {
	e.interest(otherAoi.Callback.(*Entity))
}

func (e *Entity) OnLeaveAOI(otherAoi *aoi.AOI) {
	e.unInterest(otherAoi.Callback.(*Entity))
}

// Destroy destroys the entity
func (e *Entity) Destroy() {
	if e.destroyed {
		return
	}
	log.Debugf("%s.Destroy ...", e)
	e.destroyEntity(false)
}

func (e *Entity) destroyEntity(isMigrate bool) {
	if !isMigrate {
		e.I.OnDestroy()
	} else {
		e.I.OnMigrateOut()
	}

	e.space.leave(e)
	e.clearRawTimers()

	if !isMigrate {
		e.Save(nil)
	}
	e.destroyed = true
	EntityManager.del(e)
}

// IsInterestedIn checks if other entity is interested by this entity
func (e *Entity) IsInterestedIn(other *Entity) bool {
	return e.interestedIn.Contains(other)
}

// DistanceTo calculates the distance between two entities
func (e *Entity) DistanceTo(other *Entity) unit.Coord {
	return e.Position.DistanceTo(other.Position)
}

func (e *Entity) String() string {
	return fmt.Sprintf("entity : %v", e.GetID())
}

// AsSpace converts entity to space (only works for space entity)
func (e *Entity) AsSpace() *Space {
	return (*Space)(unsafe.Pointer(e))
}

//本地调用
func (e *Entity) OnCallFromLocal(methodName string, args []interface{}) error {
	rpcDesc := e.desc.methodDesc[methodName]
	if rpcDesc == nil {
		return errors.New(fmt.Sprintf("%s.OnCallFromLocal: Method %s is not a valid RPC, args=%v", e, methodName, args))
	}

	// rpc call from server
	if rpcDesc.Flags&rfServer == 0 {
		return errors.New(fmt.Sprintf("%s.OnCallFromLocal: Method %s can not be called from Server: flags=%v", e, methodName, rpcDesc.Flags))
	}

	if rpcDesc.NumArgs < len(args) {
		return errors.New(fmt.Sprintf("%s.OnCallFromLocal: Method %s receives %d arguments, but given %d", e, methodName, rpcDesc.NumArgs, len(args)))
	}

	methodType := rpcDesc.MethodType
	in := make([]reflect.Value, rpcDesc.NumArgs+1)
	in[0] = e.V // first argument is the bind instance (self)

	for i, arg := range args {
		argType := methodType.In(i + 1)
		in[i+1] = util.Convert(arg, argType)
	}

	for i := len(args); i < rpcDesc.NumArgs; i++ { // use zero value for missing arguments
		argType := methodType.In(i + 1)
		in[i+1] = reflect.Zero(argType)
	}

	rpcDesc.Func.Call(in)
	return nil
}

//远程调用
func (e *Entity) onCallFromRemote(caller EntityID, methodName string, args [][]byte) error {
	rpcDesc := e.desc.methodDesc[methodName]
	if rpcDesc == nil {
		return errors.New(fmt.Sprintf("%s.onCallFromRemote: Method %s is not a valid RPC, args=%v", e, methodName, args))
	}

	methodType := rpcDesc.MethodType

	isFromOwnClient := e.GetID() == caller

	if rpcDesc.Flags&rfOwnClient == 0 && isFromOwnClient {
		return errors.New(fmt.Sprintf("%s.onCallFromRemote: Method %s can not be called from OwnClient: flags=%v", e, methodName, rpcDesc.Flags))
	} else if rpcDesc.Flags&rfOtherClient == 0 && !isFromOwnClient {
		return errors.New(fmt.Sprintf("%s.onCallFromRemote: Method %s can not be called from OtherClient: flags=%v, OwnClient=%v, OtherClient=%v", e, methodName, rpcDesc.Flags, e.GetID(), caller))
	}

	if rpcDesc.NumArgs < len(args) {
		return errors.New(fmt.Sprintf("%s.onCallFromRemote: Method %s receives %d arguments, but given %d", e, methodName, rpcDesc.NumArgs, len(args)))
	}

	in := make([]reflect.Value, rpcDesc.NumArgs+1)
	in[0] = e.V // first argument is the bind instance (self)

	for i, arg := range args {
		argType := methodType.In(i + 1)
		argValPtr := reflect.New(argType)
		err := msgpack.Unmarshal(arg, argValPtr.Interface())
		if err != nil {
			return errors.New(fmt.Sprintf("%s.onCallFromRemote: Method %s parse argument invalid : %v", e, methodName, err))
		}
		in[i+1] = reflect.Indirect(argValPtr)
	}

	for i := len(args); i < rpcDesc.NumArgs; i++ { // use zero value for missing arguments
		argType := methodType.In(i + 1)
		in[i+1] = reflect.Zero(argType)
	}
	rpcDesc.Func.Call(in)
	return nil
}

// Interests and Uninterest among entities
func (e *Entity) interest(other *Entity) {
	e.interestedIn.Add(other)
	other.interestedBy.Add(e)
}

func (e *Entity) unInterest(other *Entity) {
	e.interestedIn.Del(other)
	other.interestedBy.Del(e)
}

func (e *Entity) GetInterest() EntitySet {
	return e.interestedIn
}

// GetYaw gets entity Yaw
func (e *Entity) GetYaw() unit.Yaw {
	return e.Yaw
}

// SetYaw sets entity Yaw
func (e *Entity) SetYaw(yaw unit.Yaw) {
	e.Yaw = yaw
}

// FaceTo let entity face to another entity by setting Yaw accordingly
func (e *Entity) FaceTo(other *Entity) {
	e.FaceToPos(other.Position)
}

// FaceTo let entity face to a specified position, setting Yaw accordingly

func (e *Entity) FaceToPos(pos unit.Vector) {
	dir := pos.Sub(e.Position)
	dir.Y = 0
	e.SetYaw(dir.DirToYaw())
}

// Can override this function in custom entity type
func (e *Entity) OnInit() {

}

// OnattrsReady is called when entity's attribute is ready
//
// Can override this function in custom entity type
func (e *Entity) OnAttrsReady() {

}

func (e *Entity) OnCreated() {

}

// OnFreeze is called when entity is freezed
//
// Can override this function in custom entity type
func (e *Entity) OnFreeze() {
}

// OnDestroy is called when entity is destroying
//
// Can override this function in custom entity type
func (e *Entity) OnDestroy() {
}

// OnRestored is called when entity is restored
//
// Can override this function in custom entity type
func (e *Entity) OnRestored() {
}

// OnMigrateOut is called when entity is migrating out
//
// Can override this function in custom entity type
func (e *Entity) OnMigrateOut() {

}

// OnMigrateIn is called when entity is migrating in
//
// Can override this function in custom entity type
func (e *Entity) OnMigrateIn() {

}

func (e *Entity) OnClientConnected() {

}

func (e *Entity) OnClientDisconnected() {

}

func (e *Entity) OnEnterSpace() {

}

func (e *Entity) OnLeaveSpace(space *Space) {

}
