package verify

import (
	"github.com/go-playground/validator/v10"
	"sync"
)

var (
	once     sync.Once
	validate *validator.Validate
)

func init() {
	once.Do(func() {
		validate = validator.New()
	})
}

func Struct(r interface{}) error {
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}

func Slice[T any](s []T) error {
	for _, v := range s {
		if err := Struct(v); err != nil {
			return err
		}
	}
	return nil
}
