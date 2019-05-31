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
	"github.com/KylinHe/aliensboot-core/common/data_structures/set"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/KylinHe/aliensboot-core/mmo/unit"
	"reflect"
	"strings"
	"time"
)

func init() {
}

const (
	rfServer      = 1 << iota
	rfOwnClient   = 1 << iota
	rfOtherClient = 1 << iota
)

const (
	AttrClient    = 1 << iota //entity本身访问的属性
	AttrAllClient = 1 << iota //所有entity能够访问的属性
	AttrPersist   = 1 << iota //是否需要持久化
)

type methodDesc struct {
	Func       reflect.Value
	Flags      uint
	MethodType reflect.Type
	NumArgs    int
}

type methodDescMap map[string]*methodDesc

func (rdm methodDescMap) visit(method reflect.Method) {
	methodName := method.Name
	var flag uint
	var rpcName string
	if strings.HasSuffix(methodName, "_Client") {
		flag |= rfServer + rfOwnClient
		rpcName = methodName[:len(methodName)-7]
	} else if strings.HasSuffix(methodName, "_AllClient") {
		flag |= rfServer + rfOwnClient + rfOtherClient
		rpcName = methodName[:len(methodName)-11]
	} else {
		// server method
		flag |= rfServer
		rpcName = methodName
	}
	methodType := method.Type
	rdm[rpcName] = &methodDesc{
		Func:       method.Func,
		Flags:      flag,
		MethodType: methodType,
		NumArgs:    methodType.NumIn() - 1, // do not count the receiver
	}
}

// EntityTypeDesc is the entity type description for registering entity types
type EntityDesc struct {
	name EntityType

	useAOI bool

	aoiDistance unit.Coord

	entityType reflect.Type

	methodDesc methodDescMap

	clientAttrs set.StringSet

	allAttrs set.StringSet

	persistAttrs set.StringSet

	persistInterval time.Duration
}

func (desc *EntityDesc) IsPersistent() bool {
	return !desc.persistAttrs.IsEmpty()
}

func (desc *EntityDesc) DefineAttr(attr string, flag uint) {
	if flag&AttrClient != 0 {
		desc.clientAttrs.Add(attr)
	}
	if flag&AttrAllClient != 0 {
		desc.allAttrs.Add(attr)
	}
	if flag&AttrPersist != 0 {
		desc.persistAttrs.Add(attr)
	}
}

func (desc *EntityDesc) SetPersistInterval(duration time.Duration) {
	desc.persistInterval = duration
}

func (desc *EntityDesc) SetUseAOI(useAOI bool, aoiDistance unit.Coord) {
	if aoiDistance < 0 {
		log.Panic("aoi distance < 0")
	}

	desc.useAOI = useAOI
	desc.aoiDistance = aoiDistance
}
