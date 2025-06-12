package cache

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"goirc/db/model"
	db "goirc/model"
)

func Get(ctx context.Context, key string) (string, error) {
	q := model.New(db.DB.DB)
	row, err := q.CacheLoad(ctx, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("CacheLoad: %w", err)
	}
	return row.Value, nil
}

func Put(ctx context.Context, key string, value string) error {
	q := model.New(db.DB.DB)

	row, err := q.CacheLoad(ctx, key)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("CacheLoad: %w", err)
		}
	}
	if row.ID != 0 {
		err = q.CacheRemove(ctx, key)
		if err != nil {
			return fmt.Errorf("CacheRemove: %w", err)
		}
	}

	_, err = q.CacheStore(ctx, model.CacheStoreParams{Key: key, Value: value})
	if err != nil {
		return fmt.Errorf("CacheStore: %w", err)
	}

	return nil
}
