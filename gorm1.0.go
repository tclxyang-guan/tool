package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

var tdb *gorm.DB

func GetDB() *gorm.DB {
	return tdb
}

type Model struct {
	ID        uint    `gorm:"primary_key;comment:'数据id'" json:"id" req:"-"`
	CreatedAt string  `json:"created_at,omitempty" comment:"创建时间" req:"-"`
	UpdatedAt string  `json:"updated_at,omitempty" comment:"修改时间" req:"-"`
	DeletedAt *string `gorm:"type:varchar(30);default:null;comment:'删除时间'" json:"deleted_at" resp:"-" req:"-"`
}

// Logger :
type Logger interface {
	Infof(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
}

// MyLogger :
type MyLogger struct {
	logger *log.Logger
}
type MysqlConf struct {
	Address         string `yaml:"address"`
	DbName          string `yaml:"dbname"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Prefix          string `yaml:"prefix"`
	MaxOpenConns    int    `yaml:"maxconns"`
	MaxIdleConns    int    `yaml:"maxidleconns"`
	ConnMaxLifetime int64  `yaml:"connmaxlifetime"`
}

// 初始化数据库
func EnableMysql(conf MysqlConf) {
	var err error
	tdb, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Address,
		conf.DbName))

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	if conf.Prefix != "" {
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return conf.Prefix + "_" + defaultTableName
		}
	}
	tdb.SingularTable(true)
	tdb.LogMode(true)
	tdb.SingularTable(true)
	// 注册钩子函数
	tdb.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	tdb.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	tdb.Callback().Delete().Replace("gorm:delete", deleteCallback)
	tdb.DB().SetMaxIdleConns(conf.MaxIdleConns)
	tdb.DB().SetMaxOpenConns(conf.MaxOpenConns)
	tdb.DB().SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime * int64(time.Millisecond)))
}
func AutoMigrate(values ...interface{}) {
	tdb.AutoMigrate(values...)
}

// // 注册新建钩子在持久化之前
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Format("2006-01-02 15:04:05")
		if createTimeField, ok := scope.FieldByName("CreatedAt"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}
		if modifyTimeField, ok := scope.FieldByName("UpdatedAt"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// 注册更新钩子在持久化之前
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("UpdatedAt", time.Now().Format("2006-01-02 15:04:05"))
	}
}

// 注册删除钩子在删除之前
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedAt")

		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v %v %v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Format("2006-01-02 15:04:05")),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}

// Println :
func (l MyLogger) Println(v ...interface{}) {
	//l.logger.Infof(v...)
	logStr := fmt.Sprintln(v)
	l.logger.Panicln("%s", logStr)
}
