package data

import (
	"fs/core"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"strings"
	"time"
)

type TableSet[Table any] struct {
	// 上下文
	DbContext *DbContext
	// 表名
	tableName string
	db        *gorm.DB
	err       error
}

func (table TableSet[Table]) SetTableName(tableName string) {
	table.tableName = tableName
	if table.db == nil {
		return
	}
	table.db.Table(table.tableName)
}

//// NewTableSet 初始化表模型
//func NewTableSet[Table any](dbContext *DbContext, tableName string, po Table) TableSet[Table] {
//	return TableSet[Table]{
//		Po:        po,
//		DbContext: dbContext,
//		tableName: tableName,
//	}
//}

// 初始化Orm
func (table TableSet[Table]) data() *gorm.DB {
	if table.db == nil { // Data Source Name，参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name
		table.db, table.err = gorm.Open(table.getDriver(), &gorm.Config{})
		if table.err != nil {
			panic(table.err.Error())
		}
		table.db.Table(table.tableName)
		table.setPool()
	}
	return table.db
}

func (table TableSet[Table]) getDriver() gorm.Dialector {
	// 参考：https://gorm.cn/zh_CN/docs/connecting_to_the_database.html
	switch strings.ToLower(table.DbContext.dbConfig.DataType) {
	case "mysql":
		return mysql.Open(table.DbContext.dbConfig.ConnectionString)
	case "postgresql":
		return postgres.Open(table.DbContext.dbConfig.ConnectionString)
	case "sqlite":
		return sqlite.Open(table.DbContext.dbConfig.ConnectionString)
	case "sqlserver":
		return sqlserver.Open(table.DbContext.dbConfig.ConnectionString)
	}
	panic("无法识别数据库类型：" + table.DbContext.dbConfig.DataType)
}

func (table TableSet[Table]) setPool() {
	sqlDB, _ := table.db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	if table.DbContext.dbConfig.PoolMinSize > 0 {
		sqlDB.SetMaxIdleConns(table.DbContext.dbConfig.PoolMinSize)
	}
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	if table.DbContext.dbConfig.PoolMaxSize > 0 {
		sqlDB.SetMaxOpenConns(table.DbContext.dbConfig.PoolMaxSize)
	}
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func (table TableSet[Table]) Select(query interface{}, args ...interface{}) TableSet[Table] {
	table.data().Select(query, args)
	return table
}

func (table TableSet[Table]) Where(query interface{}, args ...interface{}) TableSet[Table] {
	table.data().Where(query)
	return table
}

func (table TableSet[Table]) Order(value interface{}) TableSet[Table] {
	table.data().Order(value)
	return table
}

func (table TableSet[Table]) Desc(fieldName string) TableSet[Table] {
	table.data().Order(fieldName + " desc")
	return table
}

func (table TableSet[Table]) Asc(fieldName string) TableSet[Table] {
	table.data().Order(fieldName + " asc")
	return table
}

func (table TableSet[Table]) ToList() []Table {
	var lst []Table
	table.data().Find(&lst)
	return lst
}

func (table TableSet[Table]) ToPageList(pageSize int, pageIndex int) core.PageList[Table] {
	offset := (pageIndex - 1) * pageSize
	var lst []Table
	table.data().Offset(offset).Limit(pageSize).Find(&lst)

	return core.NewPageList[Table](lst, table.Count())
}

func (table TableSet[Table]) ToEntity() Table {
	var entity Table
	table.data().First(&entity)
	return entity
}

func (table TableSet[Table]) Count() int64 {
	var count int64
	table.data().Count(&count)
	return count
}

func (table TableSet[Table]) IsExists() bool {
	var count int64
	table.data().Count(&count)
	return count > 0
}

func (table TableSet[Table]) Insert(po *Table) {
	table.data().Create(po)
}

func (table TableSet[Table]) Update(po Table) int64 {
	result := table.data().Updates(po)
	return result.RowsAffected
}

func (table TableSet[Table]) UpdateValue(column string, value interface{}) {
	table.data().Update(column, value)
}

func (table TableSet[Table]) Delete() int64 {
	result := table.data().Delete(nil)
	return result.RowsAffected
}

func (table TableSet[Table]) GetString(fieldName string) string {
	rows, _ := table.data().Select(fieldName).Limit(1).Rows()
	defer rows.Close()
	var val string
	for rows.Next() {
		rows.Scan(&val)
		// ScanRows 方法用于将一行记录扫描至结构体
		//table.data().ScanRows(rows, &user)
	}
	return val
}

func (table TableSet[Table]) GetInt(fieldName string) int {
	rows, _ := table.data().Select(fieldName).Limit(1).Rows()
	defer rows.Close()
	var val int
	for rows.Next() {
		rows.Scan(&val)
		// ScanRows 方法用于将一行记录扫描至结构体
		//table.data().ScanRows(rows, &user)
	}
	return val
}

func (table TableSet[Table]) GetLong(fieldName string) int64 {
	rows, _ := table.data().Select(fieldName).Limit(1).Rows()
	defer rows.Close()
	var val int64
	for rows.Next() {
		rows.Scan(&val)
		// ScanRows 方法用于将一行记录扫描至结构体
		//table.data().ScanRows(rows, &user)
	}
	return val
}

func (table TableSet[Table]) GetBool(fieldName string) bool {
	rows, _ := table.data().Select(fieldName).Limit(1).Rows()
	defer rows.Close()
	var val bool
	for rows.Next() {
		rows.Scan(&val)
		// ScanRows 方法用于将一行记录扫描至结构体
		//table.data().ScanRows(rows, &user)
	}
	return val
}

func (table TableSet[Table]) GetFloat32(fieldName string) float32 {
	rows, _ := table.data().Select(fieldName).Limit(1).Rows()
	defer rows.Close()
	var val float32
	for rows.Next() {
		rows.Scan(&val)
		// ScanRows 方法用于将一行记录扫描至结构体
		//table.data().ScanRows(rows, &user)
	}
	return val
}

func (table TableSet[Table]) GetFloat64(fieldName string) float64 {
	rows, _ := table.data().Select(fieldName).Limit(1).Rows()
	defer rows.Close()
	var val float64
	for rows.Next() {
		rows.Scan(&val)
		// ScanRows 方法用于将一行记录扫描至结构体
		//table.data().ScanRows(rows, &user)
	}
	return val
}
