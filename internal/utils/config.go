package utils

import (
	"os"
	"strconv"

	"github.com/usvacloud/usva/internal/dbengine"
)

func MustInt(f int, err error) int {
	if err != nil {
		return 0
	}
	return f
}

func NewTestDatabaseConfiguration() dbengine.DbConfig {
	return dbengine.DbConfig{
		Host:     StringOr(os.Getenv("DB_HOST"), "127.0.0.1"),
		Port:     IntOr(MustInt(strconv.Atoi(os.Getenv("DB_PORT"))), 5432),
		User:     StringOr(os.Getenv("DB_USERNAME_TESTS"), "usva_tests"),
		Password: StringOr(os.Getenv("DB_PASSWORD_TESTS"), "testrunner"),
		Name:     StringOr(os.Getenv("DB_NAME_TESTS"), "usva_tests"),
		UseSSL:   false,
	}
}
