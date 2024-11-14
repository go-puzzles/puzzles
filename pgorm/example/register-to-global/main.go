package main

import (
	"context"
	
	"github.com/go-puzzles/puzzles/pgorm"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
}

func (u *User) TableName() string {
	return "users"
}

func main() {
	conf := &pgorm.MysqlConfig{
		Instance: "localhost:3306",
		Database: "test",
		Username: "root",
		Password: "password",
	}
	
	_ = pgorm.RegisterSqlModelWithConf(conf, &User{})
	
	db := pgorm.GetDbByConf(conf)
	// pgorm.GetDbByModel(&User{})
	
	db.WithContext(context.Background())
}
