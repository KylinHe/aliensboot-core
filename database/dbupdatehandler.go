/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/7/31
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package database

//import (
//	"github.com/KylinHe/aliensboot-core/common/task"
//	"github.com/name5566/leaf/log"
//	"sync"
//	"time"
//	"github.com/KylinHe/aliensboot-core/common/util"
//	"reflect"
//)
//
////数据库操作枚举
//type DB_OP int
//
//const (
//	OP_INSERT DB_OP = iota
//	OP_UPDATE
//	OP_DELETE
//)
//
////数据库操作
//type DBOperation struct {
//	operation DB_OP
//	data      interface{}
//}
//
//func (this *DBOperation) change(op DB_OP, data interface{}) {
//	if op == OP_DELETE || this.operation != OP_INSERT {
//		this.operation = op
//	}
//	this.data = data
//}
//
////数据库更新配置
//type DBUpdateConfig struct {
//	//UpdateDuration time.Duration 	//缓存队列写入数据库的间隔时间
//	UpdateInterval int64 //缓存队列写入数据库的间隔时间 单位 秒
//	MaxCacheSize   int   //缓存队列的最大值，超出就需要开始写数据库了
//	MaxWriteSize   int   //缓存单次写的最大值，超出的留到下一次时间写数据库
//}
//
////数据库更新处理类 不提供给外部创建
//type dbUpdateHandler struct {
//	sync.RWMutex                         //安全锁
//	config       *DBUpdateConfig         //配置信息
//	handler      IDatabaseHandler        //数据库处理handler
//	updateCache  map[string]*DBOperation //更新队列 存储需要更新的数据 key data name + data id
//	task         *task.TimerTask         //缓存写入任务定时器
//	channel      chan *DBOperation       //写数据消息管道
//
//}
//
//func NewDBUpdateHandler(handler IDatabaseHandler, config *DBUpdateConfig) *dbUpdateHandler {
//	if handler == nil || config == nil {
//		return nil
//	}
//	//, fmt.Errorf("database handler or database update configuration can not be nil")
//	result := &dbUpdateHandler{handler: handler, config: config}
//	result.init()
//	return result
//}
//
////开启缓存写队列
//func (this *dbUpdateHandler) init() {
//	this.updateCache = make(map[string]*DBOperation)
//	if this.config.UpdateInterval > 0 {
//		this.task = &task.TimerTask{Ticker: time.NewTicker(time.Duration(time.Duration(this.config.UpdateInterval) * time.Second))}
//		this.task.Start(this.ExecuteQueue)
//	}
//	//开启处理写数据的消息处理管道
//	this.channel = make(chan *DBOperation, this.config.MaxWriteSize)
//	go func() {
//		for {
//			//只要消息管道没有关闭，就一直等待用户请求
//			dbOperation, open := <-this.channel
//			if !open {
//				break
//			}
//			this.handleOperation(dbOperation)
//		}
//		this.Close()
//	}()
//}
//
////把缓存剩余的数据写入数据库
//func (this *dbUpdateHandler) Close() {
//	this.Lock()
//	defer this.Unlock()
//	for _, dbOperation := range this.updateCache {
//		this.handleOperation(dbOperation)
//	}
//	this.updateCache = make(map[string]*DBOperation)
//	//this.task.Close()
//	//close(this.channel)
//}
//
//func (this *dbUpdateHandler) handleOperation(dbOperation *DBOperation) {
//	if dbOperation == nil {
//		return
//	}
//	begin := time.Now()
//	if dbOperation.operation == OP_DELETE {
//		this.handler.DeleteOne(dbOperation.data)
//	} else if dbOperation.operation == OP_INSERT {
//		this.handler.Insert(dbOperation.data)
//	} else if dbOperation.operation == OP_UPDATE {
//		this.handler.UpdateOne(dbOperation.data)
//	}
//	result := time.Now().Sub(begin).Seconds()
//	//记录数据库操作时间超过一秒的操作
//	if result > 1 {
//		log.Debug("execute opration op %v time(s) %v data %v", dbOperation.operation, result,
//			reflect.TypeOf(dbOperation.data))
//	}
//}
//
////加入缓存队列
//func (this *dbUpdateHandler) UpdateQueue(operation DB_OP, data interface{}) {
//	this.Lock()
//	defer this.Unlock()
//	key, err := this.getDataKey(data)
//	if err != nil {
//		log.Error("get data key error : ",err)
//		return
//	}
//	operationHistory := this.updateCache[key]
//	if operationHistory != nil {
//		operationHistory.change(operation, data)
//	} else {
//		this.updateCache[key] = &DBOperation{
//			operation: operation,
//			data:      data,
//		}
//	}
//	//if len(this.updateCache) >= this.config.MaxCacheSize {
//	//	this.ExecuteQueue()
//	//}
//}
//
//func (this *dbUpdateHandler)getDataKey(data interface{}) (string, error) {
//	key, err := this.handler.GetTableName(data)
//	if err != nil {
//		return "", err
//	}
//	id := this.handler.GetID(data)
//	switch id.(type) {
//	case string:
//		key = key + id.(string)
//		break
//	case int32:
//		key = key + util.Int32ToString(id.(int32))
//		break
//	}
//	return key, nil
//}
//
////实时更新
//func (this *dbUpdateHandler) UpdateQueueNow(operation DB_OP, data interface{}) {
//	this.Lock()
//	defer this.Unlock()
//	key, err := this.getDataKey(data)
//	if err != nil {
//		log.Error("get data key error : ", err)
//		return
//	}
//	operationHistory := this.updateCache[key]
//	if operationHistory != nil {
//		//只有删除操作和非插入操作才允许覆盖
//		operationHistory.change(operation, data)
//		delete(this.updateCache, key)
//	}
//	this.handleOperation(operationHistory)
//}
//
////将缓存的所有数据写入数据库
//func (this *dbUpdateHandler) ExecuteQueue() {
//	this.Lock()
//	defer this.Unlock()
//	for key, dbOperation := range this.updateCache {
//		select {
//		case this.channel <- dbOperation:
//		default:
//			//false 消息管道满了，直接返回, 丢到下一批次处理
//			return
//		}
//		delete(this.updateCache, key)
//	}
//}
