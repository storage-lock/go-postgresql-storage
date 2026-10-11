package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	postgresql_storage "github.com/storage-lock/go-postgresql-storage"
	"github.com/stretchr/testify/assert"
)

func TestNewPostgresqlStorageFromProperties_InitFailsWithoutDatabase(t *testing.T) {
	_, err := newPostgresqlStorageFromProperties(context.Background(), "127.0.0.1", 1, "postgres", "pass", "db")
	assert.Error(t, err)
}

func TestMain_PrintsStorageName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectExec("SET search_path").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))
	orig := postgresql_storage.PostgresqlOpen
	postgresql_storage.PostgresqlOpen = func(string, string) (*sql.DB, error) { return db, nil }
	defer func() { postgresql_storage.PostgresqlOpen = orig }()
	assert.NotPanics(t, main)
}

func TestMain_PanicsWhenInitFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectExec(".").WillReturnError(assert.AnError)
	orig := postgresql_storage.PostgresqlOpen
	postgresql_storage.PostgresqlOpen = func(string, string) (*sql.DB, error) { return db, nil }
	defer func() { postgresql_storage.PostgresqlOpen = orig }()
	assert.Panics(t, main)
}
