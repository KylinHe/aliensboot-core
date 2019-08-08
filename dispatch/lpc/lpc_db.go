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
	database2 "github.com/KylinHe/aliensboot-core/database"
	"github.com/KylinHe/aliensboot-core/module/database"
	"github.com/KylinHe/aliensboot-core/module/database/constant"
)

var DBServiceProxy = &dbHandler{}

type dbHandler struct {
}

func (handler *dbHandler) Insert(data database2.IData, dbHandler database2.IDatabaseHandler) {
	database.ChanRPC.Go(constant.DB_COMMAND_INSERT, data.Copy(), dbHandler)
}

func (handler *dbHandler) Update(data database2.IData, dbHandler database2.IDatabaseHandler) {
	database.ChanRPC.Go(constant.DB_COMMAND_UPDATE, data.Copy(), dbHandler)
}

func (handler *dbHandler) ForceUpdate(data database2.IData, dbHandler database2.IDatabaseHandler) {
	database.ChanRPC.Go(constant.DB_COMMAND_FUPDATE, data.Copy(), dbHandler)
}

func (handler *dbHandler) Delete(data database2.IData, dbHandler database2.IDatabaseHandler) {
	database.ChanRPC.Go(constant.DB_COMMAND_DELETE, data.Copy(), dbHandler)
}

func (handler *dbHandler) UpdateCondition(collectionName string, selectDoc interface{}, updateDoc interface{}, dbHandler database2.IDatabaseHandler) {
	database.ChanRPC.Go(constant.DB_COMMAND_CONDITION_UPDATE, collectionName, selectDoc, updateDoc, dbHandler)
}
