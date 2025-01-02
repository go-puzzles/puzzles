package pgorm

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-puzzles/puzzles/dialer"
	"github.com/go-puzzles/puzzles/dialer/mysql"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	thirdparty "github.com/go-puzzles/puzzles/plog/third-party"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

type MysqlDsn struct {
	DSN string
}

func (m *MysqlDsn) Validate() error {
	if m.DSN == "" {
		return errors.New("mysql config need DSN")
	}

	return nil
}

func (m *MysqlDsn) GetDBType() DbType {
	return DbMysql
}

func (m *MysqlDsn) GetService() string {
	return m.DSN
}

func (m *MysqlDsn) GetUid() string {
	return m.DSN
}

func (m *MysqlDsn) DialGorm(opts ...DialOptionFunc) (*gorm.DB, error) {
	dsnConf, err := mysqlDriver.ParseDSN(m.DSN)
	if err != nil {
		return nil, errors.Wrap(err, "parseMysqlDsn")
	}

	logPrefix := fmt.Sprintf("mysql:%s", dsnConf.DBName)

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

	return mysql.DialMysqlGormWithDSN(m.DSN, dialer.WithLogger(thirdparty.NewGormLogger(loggerOpt...)))
}

type MysqlConfig struct {
	Instance string
	Database string
	Username string
	Password string
}

func (m *MysqlConfig) Validate() error {
	if m.Instance == "" {
		return errors.New("mysql config need instace name")
	}
	if m.Database == "" {
		return errors.New("mysql config need database")
	}
	return nil
}

func (m *MysqlConfig) GetDBType() DbType {
	return DbMysql
}

func (m *MysqlConfig) GetService() string {
	return m.Instance
}

func (m *MysqlConfig) GetUid() string {
	return fmt.Sprintf("mysql-%v-%v", m.Instance, m.Database)
}

func (m *MysqlConfig) DialGorm(opts ...DialOptionFunc) (*gorm.DB, error) {
	m.TrimSpace()
	logPrefix := fmt.Sprintf("mysql:%s", m.Database)

	dialOpt := &DialOption{
		LogPrefix:     logPrefix,
		SlowThreshold: time.Millisecond * 200,
	}

	fmt.Println(opts)
	for _, optFunc := range opts {
		optFunc(dialOpt)
	}

	loggerOpt := []thirdparty.GormLoggerOption{
		thirdparty.WithPrefix(dialOpt.LogPrefix),
		thirdparty.WithSlowThreshold(dialOpt.SlowThreshold),
	}
	fmt.Println(dialOpt)
	if dialOpt.IgnoreRecordNotFound {
		loggerOpt = append(loggerOpt, thirdparty.WithIgnoreRecordNotFound())
	}

	return mysql.DialMysqlGorm(
		m.Instance,
		dialer.WithAuth(m.Username, m.Password),
		dialer.WithDBName(m.Database),
		dialer.WithLogger(thirdparty.NewGormLogger(loggerOpt...)),
	)
}

func (m *MysqlConfig) TrimSpace() {
	m.Username = strings.TrimSpace(m.Username)
	m.Password = strings.TrimSpace(m.Password)
	m.Instance = strings.TrimSpace(m.Instance)
	m.Database = strings.TrimSpace(m.Database)
}
