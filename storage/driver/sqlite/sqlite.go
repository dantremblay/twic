package sqlite

import (
	"path"

	"github.com/kassisol/twic/storage"
	"github.com/kassisol/twic/storage/driver"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/glebarez/sqlite"
)

func init() {
	storage.RegisterDriver("sqlite", New)
}

type Config struct {
	DB *gorm.DB
}

func New(config string) (driver.Storager, error) {
	dbFilePath := path.Join(config, "data.db")

	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Cert{}, &Profile{}); err != nil {
		return nil, err
	}

	return &Config{DB: db}, nil
}

func (c *Config) End() {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}
