package postgresql_storage

import (
	"database/sql"
	"fmt"
	"github.com/storage-lock/go-storage"
)

// PostgresqlStorageOptions 创建postgresql时的各种参数选项
type PostgresqlStorageOptions struct {

	// 存在哪个schema下，默认是public
	Schema string

	// 存放锁的表的名字，默认是storage_lock
	TableName string

	// 用于获取数据库连接
	ConnectionManager storage.ConnectionManager[*sql.DB]
}

func NewPostgresqlStorageOptions() *PostgresqlStorageOptions {
	return &PostgresqlStorageOptions{
		Schema:    DefaultPostgresqlStorageSchema,
		TableName: storage.DefaultStorageTableName,
	}
}

func (x *PostgresqlStorageOptions) SetSchema(schema string) *PostgresqlStorageOptions {
	x.Schema = schema
	return x
}

func (x *PostgresqlStorageOptions) GetSchema() string {
	if x.Schema == "" {
		x.Schema = DefaultPostgresqlStorageSchema
	}
	return x.Schema
}

func (x *PostgresqlStorageOptions) SetTableName(tableName string) *PostgresqlStorageOptions {
	x.TableName = tableName
	return x
}

func (x *PostgresqlStorageOptions) GetTableName() string {
	if x.TableName == "" {
		x.TableName = storage.DefaultStorageTableName
	}
	return x.TableName
}

func (x *PostgresqlStorageOptions) GetTableFullName() string {
	return fmt.Sprintf("%s.%s", x.GetSchema(), x.GetTableName())
}

func (x *PostgresqlStorageOptions) SetConnectionManager(connectionManager storage.ConnectionManager[*sql.DB]) *PostgresqlStorageOptions {
	x.ConnectionManager = connectionManager
	return x
}

func (x *PostgresqlStorageOptions) Check() error {

	if x.Schema == "" {
		x.Schema = DefaultPostgresqlStorageSchema
	}

	if x.TableName == "" {
		x.TableName = storage.DefaultStorageTableName
	}

	if x.ConnectionManager == nil {
		return fmt.Errorf("ConnectionManager can not nil")
	}

	return nil
}
