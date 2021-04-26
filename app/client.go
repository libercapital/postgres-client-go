package postgresql_client

import (
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const CONNECTED = "Successfully connected to database"

type dbclient struct {
	DB *gorm.DB
}

var instance *dbclient = nil
var once sync.Once
var dbError error = nil

/*
	@param dsn Data source name
	example: dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	@param config Gorm Configs
	example: &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		}
*/
func Connect(dsn string, config gorm.Config) (*dbclient, error) {
	once.Do(func() {
		db, err := gorm.Open(postgres.Open(dsn), &config)

		if err != nil {
			dbError = err
			return
		}

		instance = &dbclient{db}
		fmt.Println(CONNECTED)
	})

	return instance, dbError
}
