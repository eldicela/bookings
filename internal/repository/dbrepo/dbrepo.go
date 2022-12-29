package dbrepo

import (
	"database/sql"

	"github.com/eldicela/bookings/internal/config"
	"github.com/eldicela/bookings/internal/repository"
)

// type postgresDBRepo struct {
// 	App *config.AppConfig
// 	DB  *sql.DB
// }

// func NewPostgressRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
// 	return &postgresDBRepo{
// 		App: a,
// 		DB:  conn,
// 	}
// }

type mysqlDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewMysqlRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &mysqlDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
