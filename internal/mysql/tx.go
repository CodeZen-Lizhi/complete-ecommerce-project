package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Transaction 在同一个数据库事务中执行 fn。
func Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	if fn == nil {
		return errors.New("transaction fn is nil")
	}
	if DB == nil {
		return errors.New("mysql db is not initialized")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
