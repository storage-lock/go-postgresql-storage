package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewStorageFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectExec("SET search_path").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))
	storage, err := newPostgresqlStorageFromDB(context.Background(), db)
	assert.NoError(t, err)
	assert.NotEmpty(t, storage.GetName())
}

func TestMain_OpenErrorPanics(t *testing.T) {
	orig := sqlOpen
	sqlOpen = func(string, string) (*sql.DB, error) { return nil, assert.AnError }
	defer func() { sqlOpen = orig }()
	assert.Panics(t, main)
}

func TestMain_PrintsStorageName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectExec("SET search_path").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))
	orig := sqlOpen
	sqlOpen = func(string, string) (*sql.DB, error) { return db, nil }
	defer func() { sqlOpen = orig }()
	assert.NotPanics(t, main)
}

func TestMain_PanicsWhenInitFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectExec(".").WillReturnError(assert.AnError)
	orig := sqlOpen
	sqlOpen = func(string, string) (*sql.DB, error) { return db, nil }
	defer func() { sqlOpen = orig }()
	assert.Panics(t, main)
}
