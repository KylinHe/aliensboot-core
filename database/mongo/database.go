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

	queryLimit int

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
	if config.QueryLimit == 0 {
		config.QueryLimit = 50
	}
	this.dbName = config.Name
	this.queryLimit = config.QueryLimit

	c, err := Dial(config)
	if err != nil {
		return err
	}
	this.tableMetas = make(map[reflect.Type]*dbconfig.TableMeta)
	this.dbContext = c
	return nil
}


func (this *Database) SetErrorHandler(handler ErrorHandler)  {
	this.errorHandler = handler
}

//清除数据库
func (this *Database) DropDatabase() error {
	s := this.dbContext.Ref()
	defer this.dbContext.UnRef(s)
	err := s.DB(this.dbName).DropDatabase()
	if err != nil && this.errorHandler != nil {
		this.errorHandler(err)
	}
	return err
}

func (this *Database) DropCollections() error {
	s := this.dbContext.Ref()
	defer this.dbContext.UnRef(s)
	names, err := s.DB(this.dbName).CollectionNames()
	if err != nil && this.errorHandler != nil {
		this.errorHandler(err)
	}
	for _, name := range names {
		if name != IdStore {
			s.DB(this.dbName).C(name).DropCollection()
		}
	}
	return err
}

func (this *Database) NextSeq(tableMeta *dbconfig.TableMeta) (int64, error) {
	return this.dbContext.NextSeq(this.dbName, IdStore, tableMeta.Name)
}

func (this *Database) Ref(data database.IData, handler func(meta *dbconfig.TableMeta, db *mgo.Collection) error) error {
	tableMeta, err := this.GetTableMeta(data)
	if err != nil {
		if this.errorHandler != nil {
			this.errorHandler(err)
		}
		return err
	}
	s := this.dbContext.Ref()
	defer this.dbContext.UnRef(s)
	collection := s.DB(this.dbName).C(tableMeta.Name)
	err = handler(tableMeta, collection)
	if err != nil && this.errorHandler != nil {
		this.errorHandler(err)
	}
	return err
}



//func (this *Database) RefSession(handler func(database *mgo.Database) error) error {
//	s := this.dbContext.Ref()
//	defer this.dbContext.UnRef(s)
//	err := handler(s.DB(this.dbName))
//	if err != nil && this.errorHandler != nil {
//		this.errorHandler(err)
//	}
//	return err
//}

func (this *Database) BoolRef(data database.IData, handler func(meta *dbconfig.TableMeta, db *mgo.Collection) (bool, error)) (bool, error) {
	tableMeta, err := this.GetTableMeta(data)
	if err != nil {
		if this.errorHandler != nil {
			this.errorHandler(err)
		}
		return false, err
	}
	s := this.dbContext.Ref()
	defer this.dbContext.UnRef(s)
	collection := s.DB(this.dbName).C(tableMeta.Name)
	result, err := handler(tableMeta, collection)
	if err != nil && this.errorHandler != nil {
		this.errorHandler(err)
	}
	return result, err
}

func (this *Database) IntRef(data database.IData, handler func(meta *dbconfig.TableMeta, db *mgo.Collection) (int, error)) (int, error) {
	tableMeta, err := this.GetTableMeta(data)
	if err != nil {
		return 0, err
	}
	s := this.dbContext.Ref()
	defer this.dbContext.UnRef(s)
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
	//this.dbContext.Close()
}
