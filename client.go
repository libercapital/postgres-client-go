package postgresql_client

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

const CONNECTED = "Successfully connected to database"

type dbclient struct {
	DB *gorm.DB
}

var instance *dbclient
var once sync.Once

/*
	@param dsn Data source name
	example: dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
*/
func PostgreSQLConnect(dsn string) *dbclient {
	once.Do(func() {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		instance = &dbclient{db}

		fmt.Print(CONNECTED)
		fmt.Println(dsn)
	})

	return instance
}
