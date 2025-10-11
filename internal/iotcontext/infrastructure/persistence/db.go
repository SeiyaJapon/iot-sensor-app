package persistence

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type DB struct {
	conn *gorm.DB
}

func NewDB() *DB {
	dsn := os.Getenv("POSTGRES_DSN")
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	return &DB{conn: conn}
}

type SensorModel struct {
	ID        string `gorm:"primaryKey"`
	DeviceID  string `gorm:"index"`
	Name      string
	Type      string
	Config    []byte `gorm:"type:jsonb"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SensorReadingModel struct {
	ID        string `gorm:"primaryKey"`
	SensorID  string `gorm:"index"`
	DeviceID  string `gorm:"index"`
	Type      string
	Value     float64
	Unit      string
	Timestamp time.Time
	Meta      []byte `gorm:"type:jsonb"`
}

type DeviceModel struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
