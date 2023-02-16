package dbengine

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/romeq/usva/internal/generated/db"
)

type DbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	UseSSL   bool
}

func Init(x DbConfig) (*db.Queries, func()) {
	psqlconfig := "sslmode=require"
	if !x.UseSSL {
		psqlconfig = "sslmode=disable"
	}

	connstr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s",
		x.User, x.Password, x.Host, x.Port, x.Name, psqlconfig)

	r, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		log.Fatalln(err)
	}

	return db.New(r), r.Close
}
