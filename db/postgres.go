package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbDriver = "postgres"
	dbHost   = "localhost"
	dbPort   = 5432
	dbUser   = "newuser"
	dbPasswd = "rootroot"
	dbName   = "testdb"
	dbGlobal *sql.DB
)

func InitDBConnection() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPasswd, dbName)

	db, err := sql.Open(dbDriver, psqlInfo)
	db.SetMaxOpenConns(500)
	if err != nil {
		panic(err)
	}
	dbGlobal = db
}
func DBConn() *sql.DB {
	return dbGlobal
}

func GORMDB() *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: DBConn(),
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	return db
}
