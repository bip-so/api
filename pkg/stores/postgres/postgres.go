package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"time"

	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// OpenDatabaseClient generate a database client
func openDatabaseClient(ctx context.Context, c *configs.PGConnectionInfo) *gorm.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Name)
	// db, err := sql.Open("postgres", connStr)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	if err := db.Ping(); err != nil {
		log.Fatal(fmt.Errorf("\nFail to connect the database.\nPlease make sure the connection info is valid %#v", c))
		return nil
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	}) // Change this to ERROR or INFO
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Println("Connected to Database: PING")
	return gormDB
}

var db *gorm.DB

func InitDB() {
	gormDB := openDatabaseClient(context.Background(), configs.GetPGConfig())
	db = gormDB
}

func GetDB() *gorm.DB {
	return db
}
