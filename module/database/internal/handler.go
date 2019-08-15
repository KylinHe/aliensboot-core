package internal

import (
	"github.com/KylinHe/aliensboot-core/aliensboot"
	"github.com/KylinHe/aliensboot-core/database"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/KylinHe/aliensboot-core/module/database/constant"
	"reflect"

	//"reflect"
)

func init() {
	// 向当前模块注册客户端发送的消息处理函数 handleMessage
	skeleton.RegisterChanRPC(constant.DB_COMMAND_INSERT, handleInsert)
	skeleton.RegisterChanRPC(constant.DB_COMMAND_UPDATE, handleUpdate)
	skeleton.RegisterChanRPC(constant.DB_COMMAND_DELETE, handleDelete)
	skeleton.RegisterChanRPC(constant.DB_COMMAND_FUPDATE, forceUpdate)
	skeleton.RegisterChanRPC(constant.DB_COMMAND_CONDITION_UPDATE, conditionUpdate)
	skeleton.RegisterChanRPC(constant.DB_COMMAND_CONDITION_DELETE, conditionDelete)
	skeleton.RegisterChanRPC(constant.DB_COMMAND_INSERT_MULTI, handleInsertMulti)


}

func handleDelete(args []interface{}) {
	handler := args[1].(database.IDatabaseHandler)
	err := handler.DeleteOne(args[0].(database.IData))
	debugLog("delete", args[0], err)
}

func handleInsert(args []interface{}) {
	handler := args[1].(database.IDatabaseHandler)
	err := handler.Insert(args[0].(database.IData))
	debugLog("insert", args[0], err)
}

func handleInsertMulti(args []interface{}) {
	handler := args[1].(database.IDatabaseHandler)
	data := args[0].([]interface{})
	err := handler.InsertMulti(data)
	debugLog("insert multi", args[0], err)
}

func handleUpdate(args []interface{}) {
	handler := args[1].(database.IDatabaseHandler)
	err := handler.UpdateOne(args[0].(database.IData))
	debugLog("update", args[0], err)
}

func forceUpdate(args []interface{}) {
	handler := args[1].(database.IDatabaseHandler)
	err := handler.ForceUpdateOne(args[0].(database.IData))
	debugLog("force update", args[0], err)
}

func conditionUpdate(args []interface{}) {
	handler := args[3].(database.IDatabaseHandler)
	err := handler.Update(args[0].(string), args[1], args[2])
	debugLog("condition update", args[0], err)
}

func conditionDelete(args []interface{}) {
	handler := args[2].(database.IDatabaseHandler)
	err := handler.DeleteOneCondition(args[0].(database.IData), args[1])
	debugLog("condition delete", args[0], err)
}

func debugLog(opt string, data interface{}, err error) {
	if aliensboot.IsDebug() {
		if err != nil {
			typeName := reflect.TypeOf(data).Name()
			log.Debugf("[%v] %v err: %v", opt, typeName, err)
		}
	}
}
