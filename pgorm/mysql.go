package pgorm

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-puzzles/puzzles/dialer"
	thirdparty "github.com/go-puzzles/puzzles/plog/third-party"
	"gorm.io/gorm"
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

func (m *MysqlDsn) GetDBType() dbType {
	return mysql
}

func (m *MysqlDsn) GetService() string {
	return m.DSN
}

func (m *MysqlDsn) GetUid() string {
	return m.DSN
}

func (m *MysqlDsn) DialGorm() (*gorm.DB, error) {
	return dialer.DialMysqlGormWithDSN(m.DSN)
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

func (m *MysqlConfig) GetDBType() dbType {
	return mysql
}

func (m *MysqlConfig) GetService() string {
	return m.Instance
}

func (m *MysqlConfig) GetUid() string {
	return fmt.Sprintf("mysql-%v-%v", m.Instance, m.Database)
}

func (m *MysqlConfig) DialGorm() (*gorm.DB, error) {
	m.TrimSpace()
	logPrefix := fmt.Sprintf("mysql:%s", m.Database)

	return dialer.DialMysqlGorm(
		m.Instance,
		dialer.WithAuth(m.Username, m.Password),
		dialer.WithDBName(m.Database),
		dialer.WithLogger(
			thirdparty.NewGormLogger(
				thirdparty.WithPrefix(logPrefix),
				thirdparty.WithSlowThreshold(time.Millisecond*200),
			),
		),
	)
}

func (m *MysqlConfig) TrimSpace() {
	m.Username = strings.TrimSpace(m.Username)
	m.Password = strings.TrimSpace(m.Password)
	m.Instance = strings.TrimSpace(m.Instance)
	m.Database = strings.TrimSpace(m.Database)
}
