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

func (s *SqliteConfig) DialGorm(opts ...DialOptionFunc) (*gorm.DB, error) {
	logPrefix := fmt.Sprintf("sqlite:%s", s.DbFile)

	dialOpt := &DialOption{
		LogPrefix:     logPrefix,
		SlowThreshold: time.Millisecond * 200,
	}
	for _, optFunc := range opts {
		optFunc(dialOpt)
	}

	loggerOpt := []thirdparty.GormLoggerOption{
		thirdparty.WithPrefix(dialOpt.LogPrefix),
		thirdparty.WithSlowThreshold(dialOpt.SlowThreshold),
	}

	if dialOpt.IgnoreRecordNotFound {
		loggerOpt = append(loggerOpt, thirdparty.WithIgnoreRecordNotFound())
	}

	return sqlite.DialSqlLiteGorm(
		s.DbFile,
		dialer.WithLogger(thirdparty.NewGormLogger(loggerOpt...)),
	)
}
