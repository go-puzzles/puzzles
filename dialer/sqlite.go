package dialer

import (
	"os"
	"path/filepath"
	
	"github.com/go-puzzles/plog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getDBFilePath(name string) string {
	if filepath.IsAbs(name) {
		return name
	}
	
	e, _ := os.Executable()
	return filepath.Join(filepath.Dir(e), name)
}

func DialSqlLiteGorm(dbFile string, opts ...OptionFunc) (*gorm.DB, error) {
	dbFile = getDBFilePath(dbFile)
	plog.Debugf("Dial sqlite db file: %v", dbFile)
	opt := packDialOption(opts...)
	
	var (
		db  *gorm.DB
		err error
	)
	
	if opt.SqliteCache {
		db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), defaultGormConfig(opt))
	} else {
		db, err = gorm.Open(sqlite.Open(dbFile), defaultGormConfig(opt))
	}
	if err != nil {
		return nil, err
	}
	
	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}
	configDB(sqlDb)
	return db, nil
}
