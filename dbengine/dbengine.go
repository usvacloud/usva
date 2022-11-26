package dbengine

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/romeq/usva/db"
)

var DB *db.Queries

type DbConfig struct {
	Host        string
	Port        uint16
	User        string
	Password    string
	Name        string
	SslDisabled bool
}

func Init(x DbConfig) {
	connstr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		x.User, x.Password, x.Host, x.Port, x.Name)

	r, err := pgx.Connect(context.Background(), connstr)
	if err != nil {
		log.Fatalln(err)
	}

	DB = db.New(r)
}
