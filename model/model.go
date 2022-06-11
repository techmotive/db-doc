package model

// DbInfo 数据库基础信息
type DbInfo struct {
	Version   string
	Charset   string
	Collation string
	DBName    string
}

// DbConfig 数据库配置
type DbConfig struct {
	DocType       int // 1. online 2. offline
	Dsn           string
	DBName        string
	ShardingRegex string
}

// Column info
type Column struct {
	ColName    string
	ColType    string
	ColKey     string
	IsNullable string
	ColComment string
	ColDefault string
}

// Table info
type Table struct {
	TableName    string
	RealTableName    string
	TableComment string
	ColList      []Column
}
