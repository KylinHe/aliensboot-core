package mongo

import (
	"github.com/KylinHe/aliensboot-core/config"
	"github.com/KylinHe/aliensboot-core/database"
	"github.com/KylinHe/aliensboot-core/database/dbconfig"
	"gopkg.in/mgo.v2"
	"os"
	"reflect"
)

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
	dbSession *Session
	database  *mgo.Database
	auth      *database.Authority

	tableMetas map[reflect.Type]*dbconfig.TableMeta
}

//初始化连接数据库
func (this *Database) Init(config config.DBConfig) error {
	this.dbName = config.Name
	if config.MaxSession == 0 {
		config.MaxSession = 100
	}

	//优先使用环境变量
	address := os.Getenv("DBAddress")
	if address != "" {
		config.Address = address
	}

	c, err := Dial(config.Address, int(config.MaxSession))
	if err != nil {
		return err
	}
	this.tableMetas = make(map[reflect.Type]*dbconfig.TableMeta)
	this.dbContext = c
	this.dbSession = this.dbContext.Ref()
	this.database = this.dbSession.DB(config.Name)
	//if (this.auth != nil) {
	//	return this.database.Login(this.auth.Username, this.auth.Password)
	//}
	return nil
}


//初始化账号密码信息
//func (this *Database) auth(username string, password string) {
//	if username != "" {
//		this.auth = &database.Authority{username, password}
//	}
//}

func (this *Database) validateConnection() {
	if !this.isConnect() {
		panic("database is not connection")
	}
}

func (this *Database) isConnect() bool {
	return this.database != nil
}

func (this *Database) Close() {
	if this.dbContext == nil {
		return
	}
	if this.dbSession != nil {
		this.dbContext.UnRef(this.dbSession)
	}
	this.dbContext.Close()
}

func (this *Database) GetHandler() database.IDatabaseHandler {
	return this
}
