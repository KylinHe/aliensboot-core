package mongo

import (
	"github.com/KylinHe/aliensboot-core/common/util"
	"github.com/KylinHe/aliensboot-core/config"
	"github.com/KylinHe/aliensboot-core/database"
	"github.com/KylinHe/aliensboot-core/database/dbconfig"
	"gopkg.in/mgo.v2"
	"os"
	"reflect"
)

var ErrNotFound = mgo.ErrNotFound

type ErrorHandler func(err error)

//type DatabaseFactory struct {
//
//}

//func (this DatabaseFactory) Create() database.IDatabase {
//	//TODO 根据参数定制
//	return &Database{}
//}

type Database struct {
	dbName    string
	dbContext *DialContext
	//dbSession *Session
	//database  *mgo.Database
	auth      *database.Authority

	errorHandler ErrorHandler

	tableMetas map[reflect.Type]*dbconfig.TableMeta
}

//初始化连接数据库
func (this *Database) Init(config config.DBConfig) error {
	//优先使用环境变量
	address := os.Getenv("DBAddress")
	if address != "" {
		config.Address = address
	}
	name := os.Getenv("DBName")
	if name != "" {
		config.Name = name
	}
	sessionNum := os.Getenv("DBMaxSession")
	if sessionNum != "" {
		config.MaxSession = uint(util.StringToInt(sessionNum))
	}
	if config.MaxSession <= 0 {
		config.MaxSession = 50
		//log.Warnf("invalid sessionNum, reset to %v", sessionNum)
	}
	if config.DialTimeout <= 0 {
		config.DialTimeout = 30
	}
	if config.SyncTimeout <= 0 {
		config.SyncTimeout = 7
	}
	if config.SocketTimeout <= 0 {
		config.DialTimeout = 60
	}

	this.dbName = config.Name

	c, err := Dial(config)
	if err != nil {
		return err
	}
	this.tableMetas = make(map[reflect.Type]*dbconfig.TableMeta)
	this.dbContext = c
	//this.dbSession = this.dbContext.Ref()
	//this.database = this.dbSession.DB(config.Name)
	//if (this.auth != nil) {
	//	return this.database.Login(this.auth.Username, this.auth.Password)
	//}
	return nil
}


func (this *Database) SetErrorHandler(handler ErrorHandler)  {
	this.errorHandler = handler
}


//原生的更新语句
//TODO 需要拓展到内存映射修改，减少开发量
func (this *Database) Update(collection string, selector interface{}, update interface{}) error {
	s := this.dbContext.Ref()
	defer s.Close()
	return s.DB(this.dbName).C(collection).Update(selector, update)
}

func (this *Database) Ref(data interface{}, handler func(meta *dbconfig.TableMeta, db *mgo.Collection) error) error {
	tableMeta, err := this.GetTableMeta(data)
	if err != nil {
		if this.errorHandler != nil {
			this.errorHandler(err)
		}
		return err
	}
	s := this.dbContext.Ref()
	defer s.Close()
	collection := s.DB(this.dbName).C(tableMeta.Name)
	err = handler(tableMeta, collection)
	if err != nil && this.errorHandler != nil {
		this.errorHandler(err)
	}
	return err
}

func (this *Database) BoolRef(data interface{}, handler func(meta *dbconfig.TableMeta, db *mgo.Collection) (bool, error)) (bool, error) {
	tableMeta, err := this.GetTableMeta(data)
	if err != nil {
		if this.errorHandler != nil {
			this.errorHandler(err)
		}
		return false, err
	}
	s := this.dbContext.Ref()
	defer s.Close()
	collection := s.DB(this.dbName).C(tableMeta.Name)
	result, err := handler(tableMeta, collection)
	if err != nil && this.errorHandler != nil {
		this.errorHandler(err)
	}
	return result, err
}

func (this *Database) IntRef(data interface{}, handler func(meta *dbconfig.TableMeta, db *mgo.Collection) (int, error)) (int, error) {
	tableMeta, err := this.GetTableMeta(data)
	if err != nil {
		return 0, err
	}
	s := this.dbContext.Ref()
	defer s.Close()
	collection := s.DB(this.dbName).C(tableMeta.Name)
	result, err := handler(tableMeta, collection)
	if err != nil && this.errorHandler != nil {
		this.errorHandler(err)
	}
	return result, err
}


//初始化账号密码信息
//func (this *Database) auth(username string, password string) {
//	if username != "" {
//		this.auth = &database.Authority{username, password}
//	}
//}


func (this *Database) Close() {
	if this.dbContext == nil {
		return
	}
	//if this.dbSession != nil {
	//	this.dbContext.UnRef(this.dbSession)
	//}
	this.dbContext.Close()
}

func (this *Database) GetHandler() database.IDatabaseHandler {
	return this
}
