package postgresql_storage

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/storage-lock/go-storage"
	"github.com/stretchr/testify/assert"
)

func TestNewPostgresqlStorageOptions(t *testing.T) {
	opts := NewPostgresqlStorageOptions()
	assert.NotNil(t, opts)
	assert.Equal(t, DefaultPostgresqlStorageSchema, opts.Schema)
	assert.Equal(t, storage.DefaultStorageTableName, opts.TableName)
	assert.Nil(t, opts.ConnectionManager)
}

func TestPostgresqlStorageOptions_SetSchema(t *testing.T) {
	opts := NewPostgresqlStorageOptions()
	assert.NotNil(t, opts)
	opts.SetSchema("custom_schema")
	assert.Equal(t, "custom_schema", opts.Schema)
}

func TestPostgresqlStorageOptions_GetSchema_Empty(t *testing.T) {
	opts := &PostgresqlStorageOptions{}
	assert.Equal(t, DefaultPostgresqlStorageSchema, opts.GetSchema())
}

func TestPostgresqlStorageOptions_SetTableName(t *testing.T) {
	opts := NewPostgresqlStorageOptions()
	assert.NotNil(t, opts)
	opts.SetTableName("custom_table")
	assert.Equal(t, "custom_table", opts.TableName)
}

func TestPostgresqlStorageOptions_GetTableName_Empty(t *testing.T) {
	opts := &PostgresqlStorageOptions{}
	assert.Equal(t, storage.DefaultStorageTableName, opts.GetTableName())
}

func TestPostgresqlStorageOptions_GetTableFullName(t *testing.T) {
	opts := NewPostgresqlStorageOptions()
	assert.NotNil(t, opts)
	opts.SetTableName("custom_table")
	opts.SetSchema("custom_schema")
	assert.Equal(t, "custom_schema.custom_table", opts.GetTableFullName())
}

func TestPostgresqlStorageOptions_SetConnectionManager(t *testing.T) {
	opts := NewPostgresqlStorageOptions()
	assert.NotNil(t, opts)
	db := &sql.DB{}
	opts.SetConnectionManager(storage.NewFixedSqlDBConnectionManager(db))
	assert.NotNil(t, opts.ConnectionManager)
}

func TestPostgresqlStorageOptions_Check_NilConnectionManager(t *testing.T) {
	opts := NewPostgresqlStorageOptions()
	assert.NotNil(t, opts)
	err := opts.Check()
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("ConnectionManager can not nil"), err)
}

func TestPostgresqlStorageOptions_Check_ValidOptions(t *testing.T) {
	opts := NewPostgresqlStorageOptions()
	assert.NotNil(t, opts)
	db := &sql.DB{}
	opts.SetConnectionManager(storage.NewFixedSqlDBConnectionManager(db))
	err := opts.Check()
	assert.NoError(t, err)
}

func TestNewPostgresqlStorage_NilOptions(t *testing.T) {
	assert.Panics(t, func() {
		NewPostgresqlStorage(context.Background(), nil)
	})
}

func TestPostgresqlStorage_GetName(t *testing.T) {
	s := &PostgresqlStorage{}
	name := s.GetName()
	assert.Equal(t, PostgresqlStorageName, name)
}

func TestPostgresqlStorage_Capabilities(t *testing.T) {
	s := &PostgresqlStorage{}
	caps := s.Capabilities()
	assert.NotNil(t, caps)
	assert.Len(t, caps, 3)
	assert.Equal(t, storage.CapabilityCAS, caps[0])
	assert.Equal(t, storage.CapabilityReliableTime, caps[1])
	assert.Equal(t, storage.CapabilityWatch, caps[2])
}
