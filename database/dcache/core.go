package dcache

import "github.com/KylinHe/aliensboot-core/database"

type DataOp uint8

const (
	OpInsert DataOp = 1
	OpUpdate DataOp = 2
	OpDelete DataOp = 3
)

type DataCache interface {
	OpData(op DataOp, data database.IData) //更新数据操作
	Flush(dbHandler database.IDatabaseHandler) //刷新数据
}


type DataOperation struct {
	op DataOp
	data database.IData
}