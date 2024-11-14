package pgorm

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

type UserModel struct {
	ID   uint `gorm:"primarykey"`
	Name string
	Age  int
}

func (um *UserModel) TableName() string {
	return "user"
}

func TestDialMysqlDB(t *testing.T) {
	mysqlConf := MysqlConfig{
		Instance: "localhost:3306",
		Database: "sql_test",
		Username: "root",
		Password: "yang4869",
	}
	
	err := RegisterSqlModelWithConf(&mysqlConf, &UserModel{})
	assert.Nil(t, err)
	var resp []*UserModel
	err = GetDbByModel(&UserModel{}).Find(&resp).Error
	assert.Nil(t, err)
	
	t.Log(resp)
}
