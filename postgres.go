package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync/atomic"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

// ErrClosed is returned when service/database is closed
var ErrClosed = errors.New("service/database is closed")

// DB is an expense database service
type DB struct {
	db *sql.DB

	// closed dictates whether the service is closed or not
	closed int32
}

// New returns a new expense database service
func New(db *sql.DB) *DB {
	return &DB{
		db: db,
	}
}

// Close closes the underlying database connection
func (db *DB) Close() error {
	if db.IsClosed() {
		return ErrClosed
	}

	if err := db.db.Close(); err != nil {
		return err
	}
	atomic.StoreInt32(&db.closed, 1)
	return nil
}

// IsClosed returns true if the service is closed
func (db *DB) IsClosed() bool {
	return atomic.LoadInt32(&db.closed) == 1
}

// Create creates a new expense in the database and returns the created expense
func (db *DB) Create(ctx context.Context, e *Expense) (*Expense, error) {
	e = e.Sanitize()
	if err := e.Validate(); err != nil {
		return nil, err
	}
	return db.Save(ctx, e)
}

func (db *DB) Get(ctx context.Context, id int64) (*Expense, error) {
	query, args := sq.Select("id", "title", "amount", "note", "tags").
		From("expenses").
		Where("id = ?", id).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	row := db.db.QueryRow(query, args...)
	e := &Expense{}
	err := row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("Get: %d %w", id, ErrNoExpense)
	}
	if err != nil {
		return nil, fmt.Errorf("Get: %w", err)
	}
	return e, nil
}

// Save saves an expense to the database and returns the saved expense
func (db *DB) Save(ctx context.Context, e *Expense) (*Expense, error) {
	query, args := sq.Insert("expenses").
		Columns("title", "amount", "note", "tags").
		Values(e.Title, e.Amount, e.Note, pq.Array(e.Tags)).
		Suffix("RETURNING id, title, amount, note, tags").
		PlaceholderFormat(sq.Dollar).
		MustSql()

	row := db.db.QueryRow(query, args...)
	err := row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	if err != nil {
		return nil, fmt.Errorf("Save: %w", err)
	}
	return e, nil
}
