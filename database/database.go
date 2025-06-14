package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/secnex/certmanager/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Connection DatabaseConnection
	Database   *gorm.DB
}

type DatabaseConnection struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func (c *DatabaseConnection) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Database)
}

func NewConnectionFromEnv() *Database {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		port = 5432
	}
	return NewConnection(
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
	)
}

func NewConnection(host string, port int, user string, password string, database string) *Database {
	connection := DatabaseConnection{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
	}
	db, err := gorm.Open(postgres.Open(connection.String()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(models.Account{})

	return &Database{
		Connection: connection,
		Database:   db,
	}
}

func (d *Database) AutoMigrate(models ...interface{}) error {
	return d.Database.AutoMigrate(models...)
}
