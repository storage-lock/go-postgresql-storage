package postgresql_storage

import (
	"context"
	storage_test_helper "github.com/storage-lock/go-storage-test-helper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewPostgresqlStorage(t *testing.T) {
	// CI中的环境变量要保持一致名称
	envName := "STORAGE_LOCK_POSTGRESQL_DSN"
	dsn := os.Getenv(envName)
	assert.NotEmpty(t, dsn)
	connectionGetter := NewPostgresqlConnectionGetterFromDsn(dsn)
	s, err := NewPostgresqlStorage(context.Background(), NewPostgresqlStorageOptions().SetConnectionManager(connectionGetter))
	assert.Nil(t, err)
	storage_test_helper.TestStorage(t, s)
}
