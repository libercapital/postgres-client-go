package postgresql_client

import "fmt"

type Config struct {
	Host        string
	Port        string
	User        string
	Password    string
	Database    string
	ServiceName string
}

func (c Config) string() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		c.Host,
		c.User,
		c.Password,
		c.Database,
		c.Port,
	)
}
