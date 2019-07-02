/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/5/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package lpc

import (
	"github.com/KylinHe/aliensboot-core/common/util"
	database2 "github.com/KylinHe/aliensboot-core/database"
	"github.com/KylinHe/aliensboot-core/module/database"
	"github.com/KylinHe/aliensboot-core/module/database/constant"
)

var DBServiceProxy = &dbHandler{}

var safe = false

type dbHandler struct {
}

func (handler *dbHandler) Insert(data interface{}, dbHandler database2.IDatabaseHandler) {
	if safe {
		data = util.DeepClone(data)
	}
	database.ChanRPC.Go(constant.DB_COMMAND_INSERT, data, dbHandler)
}

func (handler *dbHandler) Update(data interface{}, dbHandler database2.IDatabaseHandler) {
	if safe {
		data = util.DeepClone(data)
	}
	database.ChanRPC.Go(constant.DB_COMMAND_UPDATE, data, dbHandler)
}

func (handler *dbHandler) ForceUpdate(data interface{}, dbHandler database2.IDatabaseHandler) {
	if safe {
		data = util.DeepClone(data)
	}
	database.ChanRPC.Go(constant.DB_COMMAND_FUPDATE, data, dbHandler)
}

func (handler *dbHandler) Delete(data interface{}, dbHandler database2.IDatabaseHandler) {
	if safe {
		data = util.DeepClone(data)
	}
	database.ChanRPC.Go(constant.DB_COMMAND_DELETE, data, dbHandler)
}

func (handler *dbHandler) UpdateCondition(collectionName string, selectDoc interface{}, updateDoc interface{}, dbHandler database2.IDatabaseHandler) {
	database.ChanRPC.Go(constant.DB_COMMAND_CONDITION_UPDATE, collectionName, selectDoc, updateDoc, dbHandler)
}
