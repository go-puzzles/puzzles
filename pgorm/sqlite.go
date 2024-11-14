package pgorm

import (
	"fmt"
	"time"
	
	"github.com/go-puzzles/puzzles/dialer"
	"github.com/go-puzzles/puzzles/dialer/sqlite"
	"gorm.io/gorm"
	
	thirdparty "github.com/go-puzzles/puzzles/plog/third-party"
)

type SqliteConfig struct {
	DbFile string
}

func (s *SqliteConfig) GetDBType() DbType {
	return DbSqlite
}

func (s *SqliteConfig) GetService() string {
	return s.DbFile
}

func (s *SqliteConfig) GetUid() string {
	return fmt.Sprintf("sqlite-%v", s.DbFile)
}

func (s *SqliteConfig) DialGorm() (*gorm.DB, error) {
	logPrefix := fmt.Sprintf("sqlite:%s", s.DbFile)
	return sqlite.DialSqlLiteGorm(
		s.DbFile,
		dialer.WithLogger(
			thirdparty.NewGormLogger(
				thirdparty.WithPrefix(logPrefix),
				thirdparty.WithSlowThreshold(time.Millisecond*200),
			),
		),
	)
}
