package pgorm

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SqlModel interface {
	TableName() string
}

type getClientFunc func(opts ...DialOptionFunc) *client

var (
	dbInstanceClientFuncMap = make(map[string]getClientFunc)
	modelDbMap              = make(map[string]string)
	models                  = make(map[string][]any)
)

func getDbByModel(m SqlModel, opts ...DialOptionFunc) *gorm.DB {
	rt := reflect.TypeOf(m)
	modelKey := fmt.Sprintf("%v-%v", rt.String(), m.TableName())

	key, exists := modelDbMap[modelKey]
	if !exists {
		panic(fmt.Sprintf("model %v has not been register", rt.String()))
	}

	clientFunc, ok := getInstanceClientFunc(key)
	if !ok {
		panic(fmt.Sprintf("db instance %v not found", key))
	}

	return clientFunc(opts...).DB()
}

func getDbByConfig(conf Config, opts ...DialOptionFunc) *gorm.DB {
	clientFunc, ok := getInstanceClientFunc(conf.GetUid())
	if !ok {
		panic(fmt.Sprintf("db instance %v not found", conf.GetUid()))
	}
	return clientFunc(opts...).DB()
}

func getInstanceClientFunc(key string) (getClientFunc, bool) {
	f, exists := dbInstanceClientFuncMap[key]
	return f, exists
}

func GetDbByModel(m SqlModel, opts ...DialOptionFunc) *gorm.DB {
	return getDbByModel(m, opts...).Model(m).WithContext(context.Background())
}

func GetDbByConf(conf Config, opts ...DialOptionFunc) *gorm.DB {
	return getDbByConfig(conf, opts...).WithContext(context.Background())
}

func registerInstance(conf Config) {
	key := conf.GetUid()
	dbInstanceClientFuncMap[key] = func() getClientFunc {
		var cli *client
		var once sync.Once

		f := func(opts ...DialOptionFunc) *client {
			once.Do(func() {
				cli = NewClient(conf, opts...)
			})
			return cli
		}

		return f
	}()
}

func RegisterSqlModelWithConf(conf Config, ms ...SqlModel) error {
	if _, exists := getInstanceClientFunc(conf.GetUid()); exists {
		return fmt.Errorf("%v has been register", conf.GetUid())
	}

	registerInstance(conf)

	for _, m := range ms {
		rt := reflect.TypeOf(m)
		modelKey := fmt.Sprintf("%v-%v", rt.String(), m.TableName())
		if key, exists := modelDbMap[modelKey]; exists {
			return fmt.Errorf("model %v has been register into %v", modelKey, key)
		}
		modelDbMap[modelKey] = conf.GetUid()

		// append SqlModel
		if models[conf.GetUid()] == nil {
			models[conf.GetUid()] = make([]any, 0)
		}
		models[conf.GetUid()] = append(models[conf.GetUid()], m)
	}
	return nil
}

// RegisterSqlModel will bind the corresponding SqlModel to the specified database connection create by Config
//
// Deprecated: This method will directly panic the error and is not recommended
func RegisterSqlModel(conf Config, ms ...SqlModel) {
	if _, exists := getInstanceClientFunc(conf.GetUid()); exists {
		panic(fmt.Sprintf("%v has been register", conf.GetUid()))
	}

	registerInstance(conf)

	for _, m := range ms {
		rt := reflect.TypeOf(m)
		modelKey := fmt.Sprintf("%v-%v", rt.String(), m.TableName())
		if key, exists := modelDbMap[modelKey]; exists {
			panic(fmt.Sprintf("model %v has been register into %v", modelKey, key))
		}
		modelDbMap[modelKey] = conf.GetUid()

		// append SqlModel
		if models[conf.GetUid()] == nil {
			models[conf.GetUid()] = make([]any, 0)
		}
		models[conf.GetUid()] = append(models[conf.GetUid()], m)
	}
}

func PrePing(conf Config, opts ...DialOptionFunc) error {
	db := GetDbByConf(conf, opts...)
	sqlDb, err := db.DB()
	if err != nil {
		return errors.Wrap(err, "getOriginDb")
	}

	plog.Debugf("%v db ping success", conf.GetUid())

	return sqlDb.Ping()
}

func AutoMigrate(conf Config, opts ...DialOptionFunc) error {
	ms, exists := models[conf.GetUid()]
	if !exists {
		panic(fmt.Sprintf("sqlConf %v not found", conf.GetUid()))
	}

	db := GetDbByConf(conf, opts...)
	if err := db.AutoMigrate(ms...); err != nil {
		return errors.Wrap(err, "AutoMigrate")
	}

	plog.Debugf("%v auto migrate success", conf.GetUid())
	return nil
}
