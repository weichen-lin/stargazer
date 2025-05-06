package infrastructure

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/weichen-lin/stargazer/db"
)

type Database struct {
	pool *pgxpool.Pool
}

// 可配置的連線池選項
type DatabaseOptions struct {
	MaxPoolSize int
	// 其他選項...
}

func NewDatabase(ctx context.Context, connStr string, opts *DatabaseOptions) (*Database, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("解析資料庫連線字串錯誤: %w", err)
	}

	if opts != nil {
		config.MaxConns = int32(opts.MaxPoolSize)
		// 設定其他選項...
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("連接資料庫錯誤: %w", err)
	}

	return &Database{pool: pool}, nil
}

func (d *Database) Run(ctx context.Context, fn func(q *db.Queries) error) error {
	conn, err := d.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("獲取資料庫連線錯誤: %w", err) // 返回錯誤
	}
	defer conn.Release()

	return fn(db.New(conn.Conn()))
}

func (d *Database) RunWithTransaction(ctx context.Context, fn func(q *db.Queries) error) error {
	conn, err := d.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("獲取資料庫連線錯誤: %w", err) // 返回錯誤
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("開啟事務錯誤: %w", err) // 返回錯誤
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				err = fmt.Errorf("事務回滾錯誤: %w (原始錯誤: %v)", rbErr, err) // 記錄回滾錯誤
			}
		}
	}()

	q := db.New(tx)
	err = fn(q)

	if err != nil {
		return err // 直接返回錯誤
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("提交事務錯誤: %w", err) // 返回錯誤
	}

	return nil
}
