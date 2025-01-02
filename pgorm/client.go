package pgorm

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DbType int

const (
	DbMysql = iota
	DbSqlite
)

type DialOption struct {
	LogPrefix            string
	IgnoreRecordNotFound bool
	SlowThreshold        time.Duration
}

type DialOptionFunc func(opt *DialOption)

func WithLogPrefix(prefix string) DialOptionFunc {
	return func(opt *DialOption) {
		opt.LogPrefix = prefix
	}
}

func WithDialIgnoreNotFound() DialOptionFunc {
	return func(opt *DialOption) {
		opt.IgnoreRecordNotFound = true
	}
}

func WithDialThreshold(threshold time.Duration) DialOptionFunc {
	return func(opt *DialOption) {
		opt.SlowThreshold = threshold
	}
}

type Config interface {
	GetDBType() DbType
	GetUid() string
	GetService() string
	DialGorm(...DialOptionFunc) (*gorm.DB, error)
}

type client struct {
	db     *gorm.DB
	config Config
}

func NewClient(conf Config, opts ...DialOptionFunc) *client {
	if conf.GetService() == "" {
		panic(fmt.Sprintf("pgorm: %v db service name can not be empty", conf.GetDBType()))
	}

	c := &client{config: conf}
	c.dial(opts...)
	return c
}

func (c *client) dial(opts ...DialOptionFunc) {
	db, err := c.config.DialGorm(opts...)
	if err != nil {
		panic(fmt.Sprintf("mqlClient: new client error: %v", err))
	}

	c.db = db
}

func (c *client) DB() *gorm.DB {
	return c.db
}
