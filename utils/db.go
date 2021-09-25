package utils

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     uint32
	DB       string
}

// Given the toml config file
func GetDB(path string) *gorm.DB {
	conf := DBConfig{}
	toml.DecodeFile(path, &conf)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.DB)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Connect to the database failed, err=" + err.Error())
	}
	return db
}
