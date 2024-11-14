package mysql

import (
	"database/sql"
	"fmt"
	
	"github.com/go-puzzles/puzzles/cores/discover"
	"github.com/go-puzzles/puzzles/dialer"
	"github.com/go-puzzles/puzzles/plog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func generateDSN(address string, opt *dialer.DialOption) string {
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

func DialMysqlGormWithDSN(dsn string, opts ...dialer.OptionFunc) (*gorm.DB, error) {
	opt := dialer.PackDialOption(opts...)
	db, err := gorm.Open(mysql.Open(dsn), dialer.DefaultGormConfig(opt))
	if err != nil {
		return nil, err
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	
	dialer.ConfigDB(sqlDB)
	return db, nil
}

func DialMysqlGorm(service string, opts ...dialer.OptionFunc) (*gorm.DB, error) {
	address := discover.GetServiceFinder().GetAddress(service)
	plog.Debugf("Discover mysql addr. Addr=%v", address)
	
	opt := dialer.PackDialOption(opts...)
	dsn := generateDSN(address, opt)
	
	return DialMysqlGormWithDSN(dsn, opts...)
}

// Deprecated: use DialMysqlGorm replace
func DialGorm(service string, opts ...dialer.OptionFunc) (*gorm.DB, error) {
	address := discover.GetServiceFinder().GetAddress(service)
	plog.Debugf("Discover mysql addr. Addr=%v", address)
	
	opt := dialer.PackDialOption(opts...)
	
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
	
	dialer.ConfigDB(sqlDB)
	return db, nil
}

func DialMysqlWithDSN(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "db ping")
	}
	
	dialer.ConfigDB(db)
	return db, nil
}

func DialMysql(service string, opts ...dialer.OptionFunc) (*sql.DB, error) {
	address := discover.GetServiceFinder().GetAddress(service)
	
	opt := dialer.PackDialOption(opts...)
	
	dsn := generateDSN(address, opt)
	return DialMysqlWithDSN(dsn)
}

func DialMysqlX(service string, opts ...dialer.OptionFunc) (*sqlx.DB, error) {
	db, err := DialMysql(service, opts...)
	if err != nil {
		return nil, err
	}
	
	dbx := sqlx.NewDb(db, "mysql")
	
	return dbx, nil
}
