package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	apperror "specommerce/orderservice/pkg/app_error"

	"github.com/uptrace/bun"
)

var ErrRecordNotFound = errors.New("not found record")

type SelectCriteria func(*bun.SelectQuery) *bun.SelectQuery

type CrudDatabaseOperation[T any] interface {
	FindAll(context.Context, ...SelectCriteria) ([]T, error)
	Get(context.Context, ...SelectCriteria) (T, error)
	Create(context.Context, T) (T, error)
	Update(context.Context, T) (T, error)
	Delete(context.Context, T) error
	CreateAll(context.Context, []T) ([]T, error)
	Exists(context.Context, ...SelectCriteria) (bool, error)
	FindById(context.Context, interface{}, ...SelectCriteria) (interface{}, error)
	DeleteById(context.Context, interface{}) (int, error)
}

type PostgresCrudDatabaseOperation[T any] struct {
	getDbFunc GetDbFunc
}

func NewPostgresCrudDatabaseOperation[T any](getDbFunc GetDbFunc) *PostgresCrudDatabaseOperation[T] {
	return &PostgresCrudDatabaseOperation[T]{getDbFunc: getDbFunc}
}

func (p *PostgresCrudDatabaseOperation[T]) FindAll(ctx context.Context, criteria ...SelectCriteria) ([]T, error) {
	var rows []T

	q := p.getDbFunc(ctx).NewSelect().Model(&rows)

	for i := range criteria {
		q.Apply(criteria[i])
	}

	err := q.Scan(ctx)
	return rows, err
}

func (p *PostgresCrudDatabaseOperation[T]) Get(ctx context.Context, criteria ...SelectCriteria) (T, error) {
	var row T

	q := p.getDbFunc(ctx).NewSelect().Model(&row)
	for i := range criteria {
		q.Apply(criteria[i])
	}

	err := q.Limit(1).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return row, ErrRecordNotFound
	}
	return row, err
}

func (p *PostgresCrudDatabaseOperation[T]) FindById(ctx context.Context, id interface{}, criteria ...SelectCriteria) (T, error) {
	var row T
	var idField string
	db := p.getDbFunc(ctx)
	q := db.NewSelect().Model(&row)
	table := db.Dialect().Tables().Get(reflect.TypeOf(row))
	if len(table.PKs) > 0 {
		idField = table.PKs[0].Name
	}
	if idField == "" {
		return row, errors.New("primary key not found")
	}
	q = q.Where(fmt.Sprintf("%s.%s = ?", table.Alias, idField), id)
	for i := range criteria {
		q.Apply(criteria[i])
	}
	err := q.Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return row, ErrRecordNotFound
	}
	return row, err
}

func (p *PostgresCrudDatabaseOperation[T]) getPrimaryKeyName(ctx context.Context, db bun.IDB, row T) (string, error) {
	var idField string
	table := db.Dialect().Tables().Get(reflect.TypeOf(row))
	if len(table.PKs) > 0 {
		idField = table.PKs[0].Name
	}
	if idField == "" {
		return "", apperror.NotFoundPrimaryKey
	}
	return idField, nil
}

func (p *PostgresCrudDatabaseOperation[T]) Exists(ctx context.Context, criteria ...SelectCriteria) (bool, error) {
	q := p.getDbFunc(ctx).NewSelect().Model((*T)(nil))
	for i := range criteria {
		q.Apply(criteria[i])
	}
	return q.Exists(ctx)
}

func (p *PostgresCrudDatabaseOperation[T]) Delete(ctx context.Context, row T) error {
	_, err := p.getDbFunc(ctx).NewDelete().Model(row).Exec(ctx)
	return err
}

func (p *PostgresCrudDatabaseOperation[T]) DeleteById(ctx context.Context, id interface{}) (int, error) {
	var row T
	errorTemplate := "failed to delete record"
	db := p.getDbFunc(ctx)
	q := db.NewDelete().Model(&row)
	idField, err := p.getPrimaryKeyName(ctx, db, row)
	if err != nil {
		return 0, fmt.Errorf(errorTemplate, err)
	}
	q = q.Where(fmt.Sprintf("%s = ?", idField), id)
	res, err := q.Exec(ctx)
	ra, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf(errorTemplate, err)
	}
	return int(ra), nil
}

func (p *PostgresCrudDatabaseOperation[T]) Create(ctx context.Context, row T) (T, error) {
	_, err := p.getDbFunc(ctx).NewInsert().Model(&row).Returning("*").Exec(ctx)
	return row, err
}

func (p *PostgresCrudDatabaseOperation[T]) Update(ctx context.Context, row T) (T, error) {
	res, err := p.getDbFunc(ctx).NewUpdate().Model(&row).WherePK().Returning("*").Exec(ctx)
	if err != nil {
		return row, errors.New("failed to update record")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return row, errors.New("failed to get rows affected")
	}
	if ra == 0 {
		return row, apperror.NotFoundIdWhenUpdate
	}
	return row, err
}

func (p *PostgresCrudDatabaseOperation[T]) CreateAll(ctx context.Context, req []T) ([]T, error) {
	_, err := p.getDbFunc(ctx).NewInsert().Model(&req).Returning("*").Exec(ctx)
	return req, err
}
