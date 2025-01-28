package database

import (
	"context"
	"fmt"
	"test-ns/config"
	"test-ns/internal/adapters/repo"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormDB struct {
	db *gorm.DB
}

type txKey struct{}

func NewGormDB(conf config.Config) *GormDB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Yakutsk", conf.DatabaseHost, conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseDB, conf.DatabasePort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &GormDB{db: db}
}

func (g GormDB) AutoMigrate(models ...interface{}) {
	g.db.AutoMigrate(models...)
}

func (g GormDB) Create(ctx context.Context, data interface{}) error {
	return g.db.WithContext(ctx).Create(data).Error
}

func (g GormDB) Update(ctx context.Context, tableName string, data, query interface{}, args ...interface{}) error {
	db := g.db
	tx := extractGormTx(ctx)
	if tx != nil {
		db = tx.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := db.WithContext(ctx).Table(tableName).Where(query, args...).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

func (g GormDB) BeginFind(ctx context.Context, value interface{}) repo.Find {
	db := g.db.Model(value)
	return &GormFind{
		db: db.WithContext(ctx),
	}
}

func (g GormDB) Delete(ctx context.Context, data interface{}, condition interface{}, args ...interface{}) error {
	return g.db.WithContext(ctx).Where(condition, args...).Delete(data).Error
}

func (f *GormFind) Clauses(clauses ...interface{}) repo.Find {
	clausesGORM := make([]clause.Expression, len(clauses))
	for i, c := range clauses {
		clausesGORM[i] = c.(clause.Expression)
	}
	f.db = f.db.Clauses(clausesGORM...)
	return f
}

func (g GormDB) Instance() interface{} {
	return g.db
}

type GormFind struct {
	db          *gorm.DB
	selectQuery interface{}
}

func (f *GormFind) Where(query interface{}, args ...interface{}) repo.Find {
	f.db = f.db.Where(query, args...)
	return f
}

func (f *GormFind) Find(result interface{}, args ...interface{}) error {
	query := f.db
	if f.selectQuery != nil {
		query = query.Select(f.selectQuery)
	}
	return query.Find(result, args...).Error
}

func (f *GormFind) First(result interface{}, args ...interface{}) error {
	query := f.db
	if f.selectQuery != nil {
		query = query.Select(f.selectQuery)
	}
	return query.First(result, args...).Error
}

func (f *GormFind) OrderBy(query string) repo.Find {
	f.db = f.db.Order(query)
	return f
}

type GormTransaction struct {
	db *gorm.DB
}

func NewGormTransaction(db repo.GSQL) *GormTransaction {
	dbGORM := db.Instance().(*gorm.DB)
	return &GormTransaction{db: dbGORM}
}

func (t *GormTransaction) Begin(ctx context.Context) context.Context {
	tx := t.db.WithContext(ctx).Begin()
	return context.WithValue(ctx, txKey{}, tx)
}

func (t *GormTransaction) Commit(ctx context.Context) error {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx.Commit().Error
	}
	return fmt.Errorf("no transaction found in context")
}

func (t *GormTransaction) Rollback(ctx context.Context) error {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx.Rollback().Error
	}
	return fmt.Errorf("no transaction found in context")
}

func extractGormTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return nil
}
