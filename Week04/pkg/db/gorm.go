package db

import (
	"fmt"

	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

type DBConfig struct {
	URL     string
	User    string
	PassWD  string
	CharSet string
	DBName  string
}

func (c *DBConfig) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", c.User, c.PassWD, c.URL, c.DBName, c.CharSet)
}

func NewGorm(config *DBConfig) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(config.String()), &gorm.Config{})
}
