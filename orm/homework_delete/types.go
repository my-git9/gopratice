package orm

import (
	"context"
	"database/sql"
)

type T any

// 用于 select 语句
type Querier interface {
	Get(ctx context.Context) (*T, error)
	// 这种设计也可以
	// Get(ctx context.Context) (T error)
	GetMulti(ctx context.Context) ([]*T, error)
}

// Executor 用于 Insert, Update, Delete
type Executor interface {
	Exec(ctx context.Context) (sql.Result, error)
}

// QueryBuilder 用于构建 Query
type QueryBuilder interface {
	Build() (*Query, error)
	// 这样也可以
	// Build() (Query, error)
}

type Query struct {
	SQL  string
	Args []any
}
