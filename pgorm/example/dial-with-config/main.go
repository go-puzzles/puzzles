package main

import (
	"context"
	
	"github.com/go-puzzles/puzzles/pgorm"
)

func main() {
	conf := &pgorm.MysqlConfig{
		Instance: "localhost:3306",
		Database: "test",
		Username: "root",
		Password: "password",
	}
	
	db, err := conf.DialGorm()
	if err != nil {
		panic(err)
	}
	
	db.WithContext(context.Background())
}
