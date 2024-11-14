package pgorm

import (
	"fmt"
	
	"gorm.io/gorm"
)

type dbType int

const (
	mysql = iota
	sqlite
)

type Config interface {
	GetDBType() dbType
	GetUid() string
	GetService() string
	DialGorm() (*gorm.DB, error)
}

type client struct {
	db     *gorm.DB
	config Config
}

func NewClient(conf Config) *client {
	if conf.GetService() == "" {
		panic(fmt.Sprintf("pgorm: %v db service name can not be empty", conf.GetDBType()))
	}
	
	c := &client{config: conf}
	c.dial()
	return c
}

func (c *client) dial() {
	db, err := c.config.DialGorm()
	if err != nil {
		panic(fmt.Sprintf("mqlClient: new client error: %v", err))
	}
	
	c.db = db
}

func (c *client) DB() *gorm.DB {
	return c.db
}
