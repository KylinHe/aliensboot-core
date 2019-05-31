package core

import (
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/vmihailenco/msgpack"
	"time"
)

func (e *Entity) setupSaveTimer() {
	e.addRawTimer(e.desc.persistInterval, e.Save)
}

func (e *Entity) IsPersistent() bool {
	return e.desc.IsPersistent()
}

// Save the entity
func (e *Entity) Save() {
	if !e.IsPersistent() {
		return
	}

	data := e.GetPersistentData()
	log.Debugf("save entity %v", data)
	EntityManager.GetHandler().Save(e.GetID(), e.GetType(), data, nil)
}

func (e *Entity) GetPersistentData() map[string]interface{} {
	return e.ToMapWithFilter(e.desc.persistAttrs.Contains)
}

func (e *Entity) GetClientData() map[string]interface{} {
	return e.ToMapWithFilter(e.desc.clientAttrs.Contains)
}

func (e *Entity) GetAllClientData() map[string]interface{} {
	return e.ToMapWithFilter(e.desc.allAttrs.Contains)
}

func (e *Entity) dumpTimers() []byte {
	if len(e.timers) == 0 {
		return nil
	}
	timers := make([]*entityTimerInfo, 0, len(e.timers))
	for _, t := range e.timers {
		timers = append(timers, t)
	}
	e.timers = nil
	data, err := msgpack.Marshal(timers)
	if err != nil {
		log.Debugf("pack timer data err : %v", err)
		return data
	}
	return data
}

func (e *Entity) restoreTimers(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	var timers []*entityTimerInfo
	if err := msgpack.Unmarshal(data, &timers); err != nil {
		return err
	}
	log.Debugf("%s: %d timers restored: %v", e, len(timers), timers)
	now := time.Now()
	for _, timer := range timers {
		tid := e.genTimerId()
		e.timers[tid] = timer
		timer.rawTimer = e.addRawCallback(timer.FireTime.Sub(now), func() {
			e.triggerTimer(tid, false)
		})
	}
	return nil
}

// GetMigrateData gets the migration data
func (e *Entity) GetMigrateData() *entityMigrateData {
	md := &entityMigrateData{
		Type:      e.GetType(),
		Attrs:     e.ToMap(), //all Attrs are migrated, without filter
		Pos:       e.Position,
		Yaw:       e.Yaw,
		TimerData: e.dumpTimers(),
	}
	return md
}

func (e *Entity) getAttrFlag(attrName string) (flag attrFlag) {
	if e.desc.allAttrs.Contains(attrName) {
		flag = afAllClient
	} else if e.desc.clientAttrs.Contains(attrName) {
		flag = afClient
	}

	return
}

//
//func (e *Entity) sendMapAttrChangeToClients(ma *MapAttr, key string, val interface{}) {
//	var flag attrFlag
//	if ma == e.attrs {
//		// this is the root attr
//		flag = e.getAttrFlag(key)
//	} else {
//		flag = ma.flag
//	}
//
//	if flag&afAllClient != 0 {
//		path := ma.getPathFromOwner()
//		e.client.sendNotifyMapAttrChange(e.GetID(), path, key, val)
//		for neighbor := range e.interestedBy {
//			neighbor.client.sendNotifyMapAttrChange(e.GetID(), path, key, val)
//		}
//	} else if flag&afClient != 0 {
//		path := ma.getPathFromOwner()
//		e.client.sendNotifyMapAttrChange(e.GetID, path, key, val)
//	}
//}
//
//func (e *Entity) sendMapAttrDelToClients(ma *MapAttr, key string) {
//	var flag attrFlag
//	if ma == e.attrs {
//		// this is the root attr
//		flag = e.getAttrFlag(key)
//	} else {
//		flag = ma.flag
//	}
//
//	if flag&afAllClient != 0 {
//		path := ma.getPathFromOwner()
//		e.client.sendNotifyMapAttrDel(e.GetID, path, key)
//		for neighbor := range e.interestedBy {
//			neighbor.client.sendNotifyMapAttrDel(e.GetID, path, key)
//		}
//	} else if flag&afClient != 0 {
//		path := ma.getPathFromOwner()
//		e.client.sendNotifyMapAttrDel(e.GetID, path, key)
//	}
//}
//
//func (e *Entity) sendMapAttrClearToClients(ma *MapAttr) {
//	if ma == e.attrs {
//		// this is the root attr
//		gwlog.Panicf("outmost e.Attrs can not be cleared")
//	}
//	flag := ma.flag
//
//	if flag&afAllClient != 0 {
//		path := ma.getPathFromOwner()
//		e.client.sendNotifyMapAttrClear(e.GetID, path)
//		for neighbor := range e.interestedBy {
//			neighbor.client.sendNotifyMapAttrClear(e.GetID, path)
//		}
//	} else if flag&afClient != 0 {
//		path := ma.getPathFromOwner()
//		e.client.sendNotifyMapAttrClear(e.GetID, path)
//	}
//}
//
//func (e *Entity) sendListAttrChangeToClients(la *ListAttr, index int, val interface{}) {
//	flag := la.flag
//
//	if flag&afAllClient != 0 {
//		// TODO: only pack 1 packet, do not marshal multiple times
//		path := la.getPathFromOwner()
//		e.client.sendNotifyListAttrChange(e.GetID, path, uint32(index), val)
//		for neighbor := range e.interestedBy {
//			neighbor.client.sendNotifyListAttrChange(e.GetID, path, uint32(index), val)
//		}
//	} else if flag&afClient != 0 {
//		path := la.getPathFromOwner()
//		e.client.sendNotifyListAttrChange(e.GetID, path, uint32(index), val)
//	}
//}
//
//func (e *Entity) sendListAttrPopToClients(la *ListAttr) {
//	flag := la.flag
//	if flag&afAllClient != 0 {
//		path := la.getPathFromOwner()
//		e.client.sendNotifyListAttrPop(e.GetID, path)
//		for neighbor := range e.interestedBy {
//			neighbor.client.sendNotifyListAttrPop(e.GetID, path)
//		}
//	} else if flag&afClient != 0 {
//		path := la.getPathFromOwner()
//		e.client.sendNotifyListAttrPop(e.GetID, path)
//	}
//}
//
//func (e *Entity) sendListAttrAppendToClients(la *ListAttr, val interface{}) {
//	flag := la.flag
//	if flag&afAllClient != 0 {
//		path := la.getPathFromOwner()
//		e.client.sendNotifyListAttrAppend(e.GetID, path, val)
//		for neighbor := range e.interestedBy {
//			neighbor.client.sendNotifyListAttrAppend(e.GetID, path, val)
//		}
//	} else if flag&afClient != 0 {
//		path := la.getPathFromOwner()
//		e.client.sendNotifyListAttrAppend(e.GetID, path, val)
//	}
//}
