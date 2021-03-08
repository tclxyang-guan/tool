/**
* @Auther:gy
* @Date:2020/10/23 11:01
 */

package tool

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"reflect"
	"time"
)

func EnableMysql2(conf MysqlConf) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	/*newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Silent, // Log level
			Colorful:      false,         // 禁用彩色打印
		},
	)
	newLogger.LogMode(logger.Silent)*/
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Address,
		conf.DbName)
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{
		//SkipDefaultTransaction: true,为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）。如果没有这方面的要求，您可以在初始化时禁用它。
		NamingStrategy: schema.NamingStrategy{ //GORM 允许用户通过覆盖默认的命名策略更改默认的命名约定，这需要实现接口 Namer
			TablePrefix:   conf.Prefix, // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true,        // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger:                                   logger.Default.LogMode(logger.Info), //允许通过覆盖此选项更改 GORM 的默认 logger
		DisableForeignKeyConstraintWhenMigrating: true,                                //注意 AutoMigrate 会自动创建数据库外键约束，您可以在初始化时禁用此功能
	})
	if err != nil {
		return nil, err
	}
	//自己定义的回调方法Register名字随意
	db.Callback().Create().Before("gorm:create").Register("gorm:update_time_stamp", updateTimeStampForCreateCallback2)
	db.Callback().Update().Before("gorm:update").Register("gorm:update_time_stamp", updateTimeStampForUpdateCallback2)
	//替换删除方法
	db.Callback().Delete().Replace("gorm:delete", deleteCallback2)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}

type Model2 struct {
	ID        uint    `gorm:"primary_key;comment:'数据id'" json:"id" req:"-"`
	CreatedAt string  `json:"created_at" comment:"创建时间" req:"-"`
	UpdatedAt string  `json:"updated_at" comment:"修改时间" req:"-"`
	DeletedAt Deleted `json:"deleted_at" gorm:"index" comment:"删除时间" req:"-" resp:"-"`
}

type Deleted sql.NullString

// Scan implements the Scanner interface.
func (n *Deleted) Scan(value interface{}) error {
	return (*sql.NullString)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n Deleted) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.String, nil
}

func (n Deleted) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String)
}

func (n *Deleted) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &n.String)
	if err == nil {
		n.Valid = true
	}
	return err
}

func (Deleted) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{SoftDeleteQueryClause{Field: f}}
}

type SoftDeleteQueryClause struct {
	Field *schema.Field
}

func (sd SoftDeleteQueryClause) Name() string {
	return ""
}

func (sd SoftDeleteQueryClause) Build(clause.Builder) {
}

func (sd SoftDeleteQueryClause) MergeClause(*clause.Clause) {
}

func (sd SoftDeleteQueryClause) ModifyStatement(stmt *gorm.Statement) {
	if _, ok := stmt.Clauses["soft_delete_enabled"]; !ok {
		if c, ok := stmt.Clauses["WHERE"]; ok {
			if where, ok := c.Expression.(clause.Where); ok && len(where.Exprs) > 1 {
				for _, expr := range where.Exprs {
					if orCond, ok := expr.(clause.OrConditions); ok && len(orCond.Exprs) == 1 {
						where.Exprs = []clause.Expression{clause.And(where.Exprs...)}
						c.Expression = where
						stmt.Clauses["WHERE"] = c
						break
					}
				}
			}
		}

		stmt.AddClause(clause.Where{Exprs: []clause.Expression{
			clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Value: nil},
		}})
		stmt.Clauses["soft_delete_enabled"] = clause.Clause{}
	}
}

func (Deleted) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{SoftDeleteDeleteClause{Field: f}}
}

type SoftDeleteDeleteClause struct {
	Field *schema.Field
}

func (sd SoftDeleteDeleteClause) Name() string {
	return ""
}

func (sd SoftDeleteDeleteClause) Build(clause.Builder) {
}

func (sd SoftDeleteDeleteClause) MergeClause(*clause.Clause) {
}

func (sd SoftDeleteDeleteClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.String() == "" {
		stmt.AddClause(clause.Set{{Column: clause.Column{Name: sd.Field.DBName}, Value: stmt.DB.NowFunc()}})

		if stmt.Schema != nil {
			_, queryValues := schema.GetIdentityFieldValuesMap(stmt.ReflectValue, stmt.Schema.PrimaryFields)
			column, values := schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

			if len(values) > 0 {
				stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
			}

			if stmt.ReflectValue.CanAddr() && stmt.Dest != stmt.Model && stmt.Model != nil {
				_, queryValues = schema.GetIdentityFieldValuesMap(reflect.ValueOf(stmt.Model), stmt.Schema.PrimaryFields)
				column, values = schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

				if len(values) > 0 {
					stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
				}
			}
		}

		stmt.AddClauseIfNotExists(clause.Update{})
		stmt.Build("UPDATE", "SET", "WHERE")
	}
}
func updateTimeStampForCreateCallback2(db *gorm.DB) {
	if db.Statement.Schema != nil {
		currentTime := getCurrentTime()
		db.Statement.SetColumn("CreatedAt", currentTime)
		db.Statement.SetColumn("UpdatedAt", currentTime)
	}
}

func updateTimeStampForUpdateCallback2(db *gorm.DB) {
	// if _, ok := db.Statement.Settings.Load("gorm:update_time_stamp"); ok {
	if db.Statement.Schema != nil {
		currentTime := getCurrentTime()
		db.Statement.SetColumn("UpdatedAt", currentTime)
		db.Statement.AddClause(clause.Where{
			Exprs: []clause.Expression{clause.Eq{Column: "deleted_at"}},
		})
	}
}
func deleteCallback2(db *gorm.DB) {
	if db.Error == nil {
		if db.Statement.Schema != nil && !db.Statement.Unscoped {
			for _, c := range db.Statement.Schema.DeleteClauses {
				db.Statement.AddClause(c)
			}
		}

		if db.Statement.SQL.String() == "" {
			db.Statement.SQL.Grow(100)
			db.Statement.AddClauseIfNotExists(clause.Delete{})

			if db.Statement.Schema != nil {
				_, queryValues := schema.GetIdentityFieldValuesMap(db.Statement.ReflectValue, db.Statement.Schema.PrimaryFields)
				column, values := schema.ToQueryValues(db.Statement.Table, db.Statement.Schema.PrimaryFieldDBNames, queryValues)

				if len(values) > 0 {
					db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
				}

				if db.Statement.ReflectValue.CanAddr() && db.Statement.Dest != db.Statement.Model && db.Statement.Model != nil {
					_, queryValues = schema.GetIdentityFieldValuesMap(reflect.ValueOf(db.Statement.Model), db.Statement.Schema.PrimaryFields)
					column, values = schema.ToQueryValues(db.Statement.Table, db.Statement.Schema.PrimaryFieldDBNames, queryValues)

					if len(values) > 0 {
						db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
					}
				}
			}

			db.Statement.AddClauseIfNotExists(clause.From{})
			db.Statement.Build("DELETE", "FROM", "WHERE")
		}

		if _, ok := db.Statement.Clauses["WHERE"]; !db.AllowGlobalUpdate && !ok {
			db.AddError(gorm.ErrMissingWhereClause)
			return
		}

		if !db.DryRun && db.Error == nil {
			//可通过输出Vars发现参数为2021-03-04 15:56:06.256这种带有尾缀的 因此重置
			db.Statement.Vars[0] = time.Now().Format("2006-01-02 15:04:05")
			result, err := db.Statement.ConnPool.ExecContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)

			if err == nil {
				db.RowsAffected, _ = result.RowsAffected()
			} else {
				db.AddError(err)
			}
		}
	}
}

func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
