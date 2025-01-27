package main

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
)

var acquire func() (*DB, func())

func TestMain(m *testing.M) {
	dbURL := GetEnv("DATABASE_URL", "postgres://postgres:postgres@localhost/postgres?sslmode=disable")
	if err := tryToMigrate(dbURL); err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	dbs := make(chan *DB, 1)
	dbs <- New(db)
	acquire = func() (*DB, func()) {
		db := <-dbs
		return db, func() {
			// Truncate the table before releasing the connection
			if _, err := db.db.ExecContext(context.TODO(), `TRUNCATE TABLE expenses RESTART IDENTITY CASCADE`); err != nil {
				panic(err)
			}
			dbs <- db
		}
	}

	code := m.Run()
	os.Exit(code)
}

func TestDB_Create(t *testing.T) {
	db, release := acquire()
	defer release()

	ctx := context.Background()
	expense := &Expense{
		Title:  "food",
		Amount: 100,
		Note:   "dinner",
		Tags:   []string{"food", "dinner"},
	}
	expense, err := db.Create(ctx, expense)
	if err != nil {
		t.Fatal(err)
	}
	if expense.ID == 0 {
		t.Fatal("expense ID must be greater than zero")
	}
}

func TestDB_Get(t *testing.T) {
	db, release := acquire()
	defer release()

	ctx := context.Background()
	_, err := db.Get(ctx, 1)
	if err == nil || !errors.Is(err, ErrNoExpense) {
		t.Fatalf("expected %v, got %v", ErrNoExpense, err)
	}

	_, _ = db.Create(ctx, &Expense{
		Title:  "food",
		Amount: 100,
		Note:   "dinner",
		Tags:   []string{"food", "dinner"},
	})
	expense, err := db.Get(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if expense.ID != 1 {
		t.Fatal("expected expense ID to be 1")
	}
}

func TestDB_Update(t *testing.T) {
	db, release := acquire()
	defer release()

	ctx := context.Background()
	_, err := db.Update(ctx, &Expense{
		Title:  "food",
		Amount: 100,
		Note:   "dinner",
		Tags:   []string{"food", "dinner"},
	})
	if err == nil || !errors.Is(err, ErrNoExpense) {
		t.Fatalf("expected %v, got %v", ErrNoExpense, err)
	}

	_, _ = db.Create(ctx, &Expense{
		ID:     1,
		Title:  "food",
		Amount: 100,
		Note:   "dinner",
		Tags:   []string{"food", "dinner"},
	})
	expense, err := db.Update(ctx, &Expense{
		ID:     1,
		Title:  "food",
		Amount: 100,
		Note:   "lunch",
		Tags:   []string{"food", "lunch"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if expense.Note != "lunch" {
		t.Fatal("expected expense note to be lunch")
	}
}

func TestDB_List(t *testing.T) {
	db, release := acquire()
	defer release()

	ctx := context.Background()
	expenses, err := db.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(expenses) != 0 {
		t.Fatal("expected no expenses")
	}

	_, _ = db.Create(ctx, &Expense{
		Title:  "food",
		Amount: 100,
		Note:   "dinner",
		Tags:   []string{"food", "dinner"},
	})
	_, _ = db.Create(ctx, &Expense{
		Title:  "drink",
		Amount: 100,
		Note:   "lunch",
		Tags:   []string{"drink", "lunch"},
	})

	expenses, err = db.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(expenses) != 2 {
		t.Fatal("expected two expenses")
	}
}
