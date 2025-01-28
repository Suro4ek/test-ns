package repo

import (
	"context"
)

type GSQL interface {
	AutoMigrate(models ...interface{})
	Create(ctx context.Context, data interface{}) error
	Update(ctx context.Context, tableName string, data, query interface{}, args ...interface{}) error
	BeginFind(ctx context.Context, value interface{}) Find
	Delete(ctx context.Context, data interface{}, condition interface{}, args ...interface{}) error
	Instance() interface{}
}
type Find interface {
	Where(query interface{}, args ...interface{}) Find
	Clauses(clauses ...interface{}) Find
	Find(result interface{}, args ...interface{}) error
	First(result interface{}, args ...interface{}) error
	OrderBy(query string) Find
}
type GTransaction interface {
	Begin(ctx context.Context) context.Context
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
