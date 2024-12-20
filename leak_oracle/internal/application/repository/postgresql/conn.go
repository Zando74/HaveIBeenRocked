package postgresql

import (
	"fmt"
	"sync"

	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/config"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	log = logger.LoggerSingleton.GetInstance()
)

func NewDBConn() *gorm.DB {
	cfg := config.ConfigSingleton.GetInstance()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)

	log.Info(logger.DatabaseConnectedMessage)

	return db

}

type DBConnection struct {
	once     sync.Once
	instance *gorm.DB
}

func (conn *DBConnection) GetInstance() *gorm.DB {
	conn.once.Do(func() {
		conn.instance = NewDBConn()
	})

	return conn.instance
}

func (conn *DBConnection) Close() {
	sqlDB, err := conn.instance.DB()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(logger.DatabaseClosedMessage)
	sqlDB.Close()
}

var DBConnectionSingleton DBConnection
