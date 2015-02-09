package g

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var DB *sql.DB

func InitDbConnPool() {
	var err error
	dbDsn := Config().DB.Dsn

	DB, err = sql.Open("mysql", dbDsn)
	if err != nil {
		log.Fatalf("sql.Open %s fail: %s", dbDsn, err)
	}

	DB.SetMaxIdleConns(Config().DB.MaxIdle)

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Ping() fail: %s", err)
	}
}
