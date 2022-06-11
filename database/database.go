package database

import (
	"database/sql"
	"fmt"
	"regexp"

	"db-doc/doc"
	"db-doc/model"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var dbConfig model.DbConfig

// Generate generate doc
func Generate(config *model.DbConfig) {
	dbConfig = *config
	db := initDB()
	defer db.Close()

	dbInfo := getDbInfo(db)
	dbInfo.DBName = config.DBName

	tables := getTableInfo(db, config.ShardingRegex)

	// create
	doc.CreateDoc(dbInfo, config.DocType, tables)
}

// InitDB 初始化数据库
func initDB() *sql.DB {
	db, err := sql.Open("mysql", dbConfig.Dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return db
}

// getDbInfo 获取数据库的基本信息
func getDbInfo(db *sql.DB) model.DbInfo {
	var (
		info       model.DbInfo
		rows       *sql.Rows
		err        error
		key, value string
	)
	// 数据库版本
	rows, err = db.Query("select @@version;")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&value)
	}
	info.Version = value
	// 字符集
	rows, err = db.Query("show variables like '%character_set_server%';")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&key, &value)
	}
	info.Charset = value
	// 排序规则
	rows, err = db.Query("show variables like 'collation_server%';")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&key, &value)
	}
	info.Collation = value
	return info
}

// getTableInfo 获取表信息
func getTableInfo(db *sql.DB, shardingRegex string) []model.Table {
	// find all tables
	tables := make([]model.Table, 0)
	rows, err := db.Query(getTableSQL())
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var table model.Table
	for rows.Next() {
		table.TableComment = ""
		rows.Scan(&table.TableName, &table.TableComment)
		if len(table.TableComment) == 0 {
			table.TableComment = table.TableName
		}
		table.RealTableName = table.TableName
		tables = append(tables, table)
	}
	tables = uniqTablesForSharding(tables, shardingRegex)
	for i := range tables {
		columns := getColumnInfo(db, tables[i].RealTableName)
		tables[i].ColList = columns
	}
	return tables
}
func uniqTablesForSharding(tableList []model.Table, shardingRegex string) (uniqTableList []model.Table) {
	re, err := regexp.Compile(shardingRegex)
	if err != nil {
		panic(fmt.Errorf("sharding regex is not correct."))
	}
	tableMap := map[string]model.Table{}

	for _, table := range tableList {
		tableName := re.ReplaceAllString(table.RealTableName, "")
		tableMap[tableName] = table
		table.TableName = tableName
	}

	for _, table := range tableMap {
		uniqTableList = append(uniqTableList, table)
	}
	return uniqTableList
}

// getColumnInfo 获取列信息
func getColumnInfo(db *sql.DB, tableName string) []model.Column {
	columns := make([]model.Column, 0)
	rows, err := db.Query(getColumnSQL(tableName))
	if err != nil {
		fmt.Println(err)
	}
	var column model.Column
	for rows.Next() {
		rows.Scan(&column.ColName, &column.ColType, &column.ColKey, &column.IsNullable, &column.ColComment, &column.ColDefault)
		columns = append(columns, column)
	}
	return columns
}

// getTableSQL
func getTableSQL() string {
	return fmt.Sprintf(`
			select table_name    as TableName, 
			       table_comment as TableComment
			from information_schema.tables 
			where table_schema = '%s'
		`, dbConfig.DBName)
}

// getColumnSQL
func getColumnSQL(tableName string) string {
	return fmt.Sprintf(`
			select column_name as ColName,
			column_type        as ColType,
			column_key         as ColKey,
			is_nullable        as IsNullable,
			column_comment     as ColComment,
			column_default     as ColDefault
			from information_schema.columns 
			where table_schema = '%s' and table_name = '%s'
		`, dbConfig.DBName, tableName)
}
