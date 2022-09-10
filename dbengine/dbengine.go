package dbengine

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/romeq/usva/utils"
)

var DbConnection *sqlx.DB
var schema string = `
CREATE TABLE IF NOT EXISTS files(
	id 			INTEGER PRIMARY KEY AUTOINCREMENT,
	filename 	VARCHAR(256) 	NOT NULL UNIQUE,
	password 	VARCHAR(512),
	upload_date VARCHAR(256) 	NOT NULL,
	file_size 	INTEGER 		NOT NULL,
	viewcount	INTEGER			NOT NULL,
	owner_id 	INTEGER
);
`

type File struct {
	ID         int
	Filename   string
	Password   string
	UploadDate string `db:"upload_date"`
	FileSize   int    `db:"file_size"`
	OwnerId    int    `db:"owner_id"`
	ViewCount  int    `db:"viewcount"`
}

func Init(datasource string) {
	// connect database
	var err error
	DbConnection, err = sqlx.Connect("sqlite3", datasource)
	utils.Check(err)

	// establishing connection and call
	utils.Check(DbConnection.Ping())
	migrate()
}

func migrate() {
	if DbConnection == nil {
		log.Fatalln("migrate function called while db is not connected")
	}

	DbConnection.MustExec(schema)
}
