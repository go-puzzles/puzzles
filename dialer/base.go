package dialer

import (
	"database/sql"
	"time"
	
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	DefaultUserName = "root"
)

type DialOption struct {
	User     string
	Password string
	DBName   string
	
	SqliteCache bool
	Logger      logger.Interface
}

type OptionFunc func(*DialOption)

func WithAuth(user, pwd string) OptionFunc {
	return func(do *DialOption) {
		do.User = user
		do.Password = pwd
	}
}

func WithDBName(db string) OptionFunc {
	return func(do *DialOption) {
		do.DBName = db
	}
}

func WithSqliteCache() OptionFunc {
	return func(do *DialOption) {
		do.SqliteCache = true
	}
}

func WithLogger(log logger.Interface) OptionFunc {
	return func(do *DialOption) {
		do.Logger = log
	}
}

func PackDialOption(opts ...OptionFunc) *DialOption {
	opt := &DialOption{}
	for _, o := range opts {
		o(opt)
	}
	
	if opt.User == "" {
		opt.User = DefaultUserName
	}
	
	return opt
}

func ConfigDB(sqlDB *sql.DB) {
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func DefaultGormConfig(opt *DialOption) *gorm.Config {
	return &gorm.Config{
		PrepareStmt: true,
		Logger:      opt.Logger,
	}
}
