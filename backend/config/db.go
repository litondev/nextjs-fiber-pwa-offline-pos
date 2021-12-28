package config

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Database(env map[string]string) (*gorm.DB, error) {
	dsn := "host=" + env["DB_HOST"] +
		" user=" + env["DB_USER"] +
		" password=" + env["DB_PASSWORD"] +
		" dbname=" + env["DB_NAME"] +
		" port=" + env["DB_PORT"] +
		" sslmode=disable TimeZone=Asia/Jakarta"

	db, errDb := gorm.Open(postgres.New(postgres.Config{
		PreferSimpleProtocol: true,
		DSN:                  dsn,
	}))

	if errDb != nil {
		// Print Error
		fmt.Println(errDb)

		// Throw Error
		return nil, errors.New("Can't Connect To Database")
	}

	// Return Gorm Orm
	return db, nil
}
