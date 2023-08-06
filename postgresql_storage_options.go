package postgresql_storage

import (
	"database/sql"
	"github.com/storage-lock/go-storage"
)

type PostgreSQLStorageOptions struct {

	// 存在哪个schema下，默认是public
	Schema string

	// 存放锁的表的名字
	TableName string

	// 用于获取数据库连接
	ConnectionManager storage.ConnectionManager[*sql.DB]
}

func NewPostgreSQLStorageOptions() *PostgreSQLStorageOptions {
	return &PostgreSQLStorageOptions{
		Schema: DefaultPostgreSQLStorageSchema,
	}
}

func (x *PostgreSQLStorageOptions) SetSchema(schema string) *PostgreSQLStorageOptions {
	x.Schema = schema
	return x
}

func (x *PostgreSQLStorageOptions) SetTableName(tableName string) *PostgreSQLStorageOptions {
	x.TableName = tableName
	return x
}

func (x *PostgreSQLStorageOptions) SetConnectionManager(connectionManager storage.ConnectionManager[*sql.DB]) *PostgreSQLStorageOptions {
	x.ConnectionManager = connectionManager
	return x
}
