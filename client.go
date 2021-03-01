package postgresql_client

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			panic(err)
		}

		instance = &dbclient{db}

		fmt.Println(CONNECTED)
	})

	return instance
}
