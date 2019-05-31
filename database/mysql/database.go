package mysql

//
//import (
//	"github.com/KylinHe/aliensboot-core/database"
//	"fmt"
//	_ "github.com/go-sql-driver/mysql"
//	"github.com/jinzhu/gorm"
//	"github.com/KylinHe/aliensboot-core/database/dbconfig"
//)
//
//type Database struct {
//	dbName string
//	db     *gorm.DB
//	*database.Authority
//}
//
//func (this *Database) getAuthorityString() string {
//	if this.Authority == nil {
//		return ""
//	}
//	return this.Username + ":" + this.Password
//}
//
////初始化账号密码信息
////func (this *Database) auth(username string, password string) {
////	if username != "" {
////		this.Authority = &database.Authority{username, password}
////	}
////}
//
////初始化连接数据库
//func (this *Database) Init(config dbconfig.DBConfig) error {
//	this.dbName = config.Name
//	db, err := gorm.Open("mysql",
//		fmt.Sprintf(this.getAuthorityString()+"@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local", config.Address, config.Name))
//	if err == nil {
//		this.db = db
//	}
//	db.LogMode(true)
//	return err
//}
//
//func (this *Database) Close() {
//	if this.db != nil {
//		this.db.Close()
//	}
//}
//
//func (this *Database) GetHandler() database.IDatabaseHandler {
//	return this
//}
