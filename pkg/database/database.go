package database

import (
	"log"

	"github.com/Doittikorn/go-e-commerce/config"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DbConnectPostgresql(cfg config.DBConfigImpl) *sqlx.DB {
	db, err := sqlx.Connect("pgx", cfg.Url())
	if err != nil {
		log.Fatalf("connect to db failed: %v", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxOpenConnection())
	return db
}

func DbConnectMySQL(cfg config.DBConfigImpl) *sqlx.DB {
	// connect to mysql

	db, err := sqlx.Connect("mysql", "root:my-secret-pw@tcp(localhost:3306)/test")
	if err != nil {
		log.Fatalf("connect to db failed: %v", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxOpenConnection())
	return db
}
