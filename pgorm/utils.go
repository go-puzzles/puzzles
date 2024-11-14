package pgorm

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func GormIsErrRecordNotFound(err error) (bool, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	} else if err != nil {
		return false, err
	} else {
		return false, nil
	}
}
