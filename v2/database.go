package postgresql_client

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres interface {
	MigrateUp(path string) error
	MigrateDown(path string) error
	Connect() error
	DB() *gorm.DB
}

type postgresClientImpl struct {
	config     config
	poolConfig poolConfig
	db         *gorm.DB
}

func (p *postgresClientImpl) Connect() error {
	gormConfig := gorm.Config{
		SkipDefaultTransaction: true,
	}

	if p.config.NamingStrategy != nil {
		gormConfig.NamingStrategy = p.config.NamingStrategy
	}

	sqltrace.Register("postgres", &pq.Driver{})

	tracedb, err := sqltrace.Open("postgres", p.config.string(), sqltrace.WithServiceName(p.config.ServiceName))

	if err != nil {
		return err
	}

	db, err := gormtrace.Open(
		gormpostgres.New(
			gormpostgres.Config{
				Conn: tracedb,
			},
		),
		&gormConfig,
		gormtrace.WithServiceName(p.config.ServiceName),
	)

	if err != nil {
		return err
	}

	p.db = db

	sqlDB, err := p.db.DB()

	if err != nil {
		return err
	}

	if p.poolConfig.MaxIdle < 1 {
		return ErrInvalidPoolConfiguration
	}

	if p.poolConfig.MaxOpen < 1 {
		return ErrInvalidPoolConfiguration
	}

	if p.poolConfig.MaxLifeTime < 1 {
		return ErrInvalidPoolConfiguration
	}

	sqlDB.SetMaxIdleConns(p.poolConfig.MaxIdle)
	sqlDB.SetMaxOpenConns(p.poolConfig.MaxOpen)
	sqlDB.SetConnMaxLifetime(p.poolConfig.MaxLifeTime)

	return nil
}

func (p *postgresClientImpl) DB() *gorm.DB {
	return p.db
}

func (p *postgresClientImpl) MigrateDown(path string) error {
	sqlInstance, err := p.db.DB()

	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(sqlInstance, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		path,
		p.config.Database,
		driver,
	)

	if err != nil {
		return err
	}

	err = m.Down()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (p *postgresClientImpl) MigrateUp(path string) error {
	sqlInstance, err := p.db.DB()

	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(sqlInstance, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		path,
		p.config.Database,
		driver,
	)

	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func New(config config, poolConfig poolConfig) Postgres {
	return &postgresClientImpl{
		config:     config,
		poolConfig: poolConfig,
	}
}
