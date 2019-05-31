package mysql

//
//import (
//	"github.com/KylinHe/aliensboot-core/database"
//)
//
////  grom关联文档  http://gorm.book.jasperxu.com/associations.html#hm
//
//var Cache struct {
//	data map[string]*DataConfig
//}
//
//type DataConfig struct {
//}
//
////
//func (this *Database) EnsureCounter(data database.IData) {
//	//this.db.Model(data).
//	//TODO 设置键自增长
//}
//
////确保含有表结构
//func (this *Database) EnsureTable(data database.IData) {
//	if !this.db.HasTable(data.Name()) {
//		this.db.Table(data.Name()).Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(data)
//	}
//}
//
////确保索引
//func (this *Database) EnsureUniqueIndex(data database.IData, key string) {
//
//}
//
////设置关联关系,这样查询能够注入
//func (this *Database) Related(data interface{}, relateData interface{}, relateTableName string, relateKey string) error {
//	return this.db.Table(relateTableName).Model(data).Related(relateData, relateKey).Error
//}
//
////自增长id
//func (this *Database) GenId(data database.IData) int32 {
//	return 0
//}
//
//func (this *Database) GenTimestampId(data database.IData) int64 {
//	return 0
//}
//
////插入新数据
//func (this *Database) Insert(data database.IData) error {
//	return this.db.Table(data.Name()).Create(data).Error
//}
//
////查询所有数据
//func (this *Database) QueryAll(data database.IData, result interface{}) error {
//	return this.db.Table(data.Name()).Find(result).Error
//}
//
////查询单条记录
//func (this *Database) QueryOne(data database.IData) error {
//	return this.db.Table(data.Name()).First(data).Error
//}
//
////func parseTagSetting(tags reflect.StructTag) map[string]string {
////	setting := map[string]string{}
////	for _, str := range []string{tags.Get("sql"), tags.Get("gorm")} {
////		tags := strings.Split(str, ";")
////		for _, value := range tags {
////			v := strings.Split(value, ":")
////			k := strings.TrimSpace(strings.ToUpper(v[0]))
////			if len(v) >= 2 {
////				setting[k] = strings.Join(v[1:], ":")
////			} else {
////				setting[k] = k
////			}
////		}
////	}
////	return setting
////}
//
//func (this *Database) DeleteOne(data database.IData) error {
//	return this.db.Table(data.Name()).Delete(data).Error
//}
//
//func (this *Database) DeleteOneCondition(data database.IData, selector interface{}) error {
//	//TODO impl
//	return nil
//}
//
////查询单条记录
//func (this *Database) IDExist(data database.IData) bool {
//	return this.db.Table(data.Name()).First(data).Error == nil
//}
//
////按条件多条查询
//func (this *Database) QueryAllCondition(data database.IData, condition string, value interface{}, result interface{}) error {
//	return this.db.Table(data.Name()).Where(condition+" = ?", value).Find(result).Error
//}
//
////按条件单条查询
//func (this *Database) QueryOneCondition(data database.IData, condition string, value interface{}) error {
//	return this.db.Table(data.Name()).Where(condition+" = ?", value).First(data).Error
//}
//
////更新单条数据
//func (this *Database) UpdateOne(data database.IData) error {
//	return this.db.Table(data.Name()).Save(data).Error
//}
//
//func (this *Database) ForceUpdateOne(data database.IData) error {
//	if this.IDExist(data) {
//		return this.UpdateOne(data)
//	} else {
//		return this.Insert(data)
//	}
//}
//
////条件更新  selector 查询语句   update 更新语句
//func (this *Database) Update(collection string, selector interface{}, update interface{}) error {
//	return this.db.Table(collection).Where(selector).Update(update).Error
//}
