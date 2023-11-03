package postgresql_client

import (
	"fmt"
	"time"

	"gorm.io/gorm/schema"
)

type poolConfig struct {
	MaxIdle     int
	MaxOpen     int
	MaxLifeTime time.Duration
}

func PoolConfig(maxIdle int, maxOpen int, maxLifetime time.Duration) poolConfig {
	return poolConfig{
		MaxIdle:     maxIdle,
		MaxOpen:     maxOpen,
		MaxLifeTime: maxLifetime,
	}
}

func Config(
	Host string,
	Port string,
	User string,
	Password string,
	Database string,
	ServiceName string,
	SSLMode string,
	NamingStrategy schema.Namer,
) config {
	return config{
		Host,
		Port,
		User,
		Password,
		Database,
		ServiceName,
		SSLMode,
		NamingStrategy,
	}
}

type config struct {
	Host           string
	Port           string
	User           string
	Password       string
	Database       string
	ServiceName    string
	SSLMode        string
	NamingStrategy schema.Namer
}

func (c config) string() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=America/Sao_Paulo",
		c.Host,
		c.User,
		c.Password,
		c.Database,
		c.Port,
		c.SSLMode,
	)
}
