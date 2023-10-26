package main

import (
	"os"

	"github.com/Doittikorn/go-e-commerce/config"
	"github.com/Doittikorn/go-e-commerce/modules/servers"
	"github.com/Doittikorn/go-e-commerce/pkg/database"
)

func envPath() string {

	if len(os.Args) == 1 {
		return ".env"
	}

	return os.Args[1]
}

func main() {

	cfg := config.LoadConfig(envPath())

	db := database.DbConnectPostgresql(cfg.DB())
	// dbMySql := database.DbConnectMySQL(cfg.DB())

	servers.NewServer(cfg, db).Start()

	defer db.Close()
}
