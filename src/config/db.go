package config

import (
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/swim233/do_not_xm/logger"
	"gorm.io/gorm"
)

var (
	ConfigDB *gorm.DB
)

type GroupInterface interface {
	GroupDB
}
type GroupDB struct {
	GroupID      int64 `gorm:"primaryKey"`
	StaticCD     time.Duration
	RandomCD     time.Duration
	Mode         string
	ListenID     string
	KeyWord      string
	ReplyMessage string
}

func (GroupDB) TableName() string {
	return "groups"
}

// 初始化数据库
func InitDB() {
	defer logger.Suger.Sync()
	if err := os.MkdirAll("db", 0755); err != nil {
		logger.Suger.Warnf("Fail to create db folder : %s", err.Error())
	}
	if db, err := gorm.Open(sqlite.Open("db/data.db"), &gorm.Config{}); err != nil {
		logger.Suger.Panicf("Fail to open database file,trying create : %s", err.Error())

	} else {
		ConfigDB = db
	}
	if err := ConfigDB.AutoMigrate(&GroupDB{}); err != nil {
		logger.Suger.Panicf("Fail to auto migrate database : %s", err.Error())
	}
}

func WriteBackDatabase(group GroupDB) {
	if result := ConfigDB.Save(&group); result.Error != nil {
		logger.Suger.Errorf("Fail to write back config : %s", result.Error.Error())
	}
}
