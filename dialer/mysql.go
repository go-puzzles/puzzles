package dialer

import (
	"database/sql"
	"fmt"
	
	"github.com/go-puzzles/cores/discover"
	"github.com/go-puzzles/plog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DefaultUserName = "root"
)

func generateDSN(address string, opt *DialOption) string {
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		opt.User,
		opt.Password,
		address,
		opt.DBName,
	)
	plog.Debugf("sql dsn generate. dsn=%v", dsn)
	return dsn
}

func DialMysqlGorm(service string, opts ...OptionFunc) (*gorm.DB, error) {
	address := discover.GetServiceFinder().GetAddress(service)
	plog.Debugf("Discover mysql addr. Addr=%v", address)
	
	opt := packDialOption(opts...)
	
	dsn := generateDSN(address, opt)
	db, err := gorm.Open(mysql.Open(dsn), defaultGormConfig(opt))
	if err != nil {
		return nil, err
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	
	configDB(sqlDB)
	return db, nil
}

// Deprecated: use DialMysqlGorm replace
func DialGorm(service string, opts ...OptionFunc) (*gorm.DB, error) {
	address := discover.GetServiceFinder().GetAddress(service)
	plog.Debugf("Discover mysql addr. Addr=%v", address)
	
	opt := packDialOption(opts...)
	
	dsn := generateDSN(address, opt)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	
	configDB(sqlDB)
	return db, nil
}

func DialMysql(service string, opts ...OptionFunc) (*sql.DB, error) {
	address := discover.GetServiceFinder().GetAddress(service)
	
	opt := packDialOption(opts...)
	
	dsn := generateDSN(address, opt)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "db ping")
	}
	
	configDB(db)
	return db, nil
}

func DialMysqlX(service string, opts ...OptionFunc) (*sqlx.DB, error) {
	db, err := DialMysql(service, opts...)
	if err != nil {
		return nil, err
	}
	
	dbx := sqlx.NewDb(db, "mysql")
	
	return dbx, nil
}
