package main

import (
	"flag"
	"fmt"
	"os"

	"db-doc/database"
	"db-doc/model"

	"github.com/go-sql-driver/mysql"
)

const version = "v1.1.1"

var dbConfig model.DbConfig
var _dsn = "sz_to_c_test:Akofm8v4NupZfsB_QoVd@(master.shopee_mkt_ap_meta.mysql.cloud.test.shopee.io:6606)/shopee_mkt_ap_payin_sg_db?charset=utf8mb4&parseTime=True&loc=Local"
var __dsn = "root:@(localhost:3306)/ferry?charset=utf8mb4&parseTime=True&loc=Local"
var dsn = flag.String("dsn","","specify the mysql dsn")
var online = flag.Bool("online",false,"gen online doc")
var shardingRegex = flag.String("shardingRegex", `_\d{8}`, "specify the sharding regex which sharding table matchs")


func main() {
	flag.Parse()

	if len(*dsn) == 0 {
		fmt.Println("need dsn")
		os.Exit(-1)
	}

	cfg, err := mysql.ParseDSN(*dsn)
	if err != nil {
		fmt.Println("invalid dsn")
		os.Exit(-1)
	}
	dbConfig.Dsn = *dsn
	dbConfig.DBName = cfg.DBName
	dbConfig.ShardingRegex = *shardingRegex
	// doc type
	dbConfig.DocType = 2
	if *online {
		dbConfig.DocType = 1
	}
	// generate
	database.Generate(&dbConfig)
}


