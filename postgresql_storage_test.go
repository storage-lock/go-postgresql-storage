package postgresql_storage

import (
	"context"
	storage_test_helper "github.com/storage-lock/go-storage-test-helper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewPostgreSQLStorage(t *testing.T) {
	envName := "STORAGE_LOCK_POSTGRESQL_DSN"
	dsn := os.Getenv(envName)
	assert.NotEmpty(t, dsn)
	connectionGetter := NewPostgreSQLConnectionGetterFromDSN(dsn)
	s, err := NewPostgreSQLStorage(context.Background(), &PostgreSQLStorageOptions{
		ConnectionManager: connectionGetter,
		TableName:         storage_test_helper.TestTableName,
	})
	assert.Nil(t, err)
	storage_test_helper.TestStorage(t, s)
}
