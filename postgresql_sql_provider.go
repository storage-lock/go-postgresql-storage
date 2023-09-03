package postgresql_storage

import (
	"context"
	"fmt"
	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
	"github.com/storage-lock/go-storage"
)

// PostgresqlSqlProvider storage sql的postgresql方言
type PostgresqlSqlProvider struct {
	*sql_based_storage.Sql92Provider
}

var _ sql_based_storage.SqlProvider = &PostgresqlSqlProvider{}

func NewPostgresqlSqlProvider() *PostgresqlSqlProvider {
	return &PostgresqlSqlProvider{}
}

func (x *PostgresqlSqlProvider) NowTimestampSql(ctx context.Context, tableFullName string) (string, []any) {
	return "SELECT CURRENT_TIMESTAMP", nil
}

func (x *PostgresqlSqlProvider) SelectLockInformationJsonStringSql(ctx context.Context, tableFullName string, lockId string) (string, []any) {
	getLockSql := fmt.Sprintf("SELECT lock_information_json_string FROM %s WHERE lock_id = $1", tableFullName)
	return getLockSql, []any{lockId}
}

func (x *PostgresqlSqlProvider) DeleteWithVersionSql(ctx context.Context, tableFullName string, lockId string, exceptedVersion storage.Version, lockInformation *storage.LockInformation) (string, []any) {
	deleteSql := fmt.Sprintf("DELETE FROM %s WHERE lock_id = $1 AND owner_id = $2 AND version = $3", tableFullName)
	return deleteSql, []any{lockId, lockInformation.OwnerId, exceptedVersion}
}

func (x *PostgresqlSqlProvider) CreateWithVersionSql(ctx context.Context, tableFullName string, lockId string, version storage.Version, lockInformation *storage.LockInformation) (string, []any) {
	insertSql := fmt.Sprintf("INSERT INTO %s (lock_id, owner_id, version, lock_information_json_string) VALUES ($1, $2, $3, $4)", tableFullName)
	return insertSql, []any{lockId, lockInformation.OwnerId, version, lockInformation.ToJsonString()}
}

func (x *PostgresqlSqlProvider) UpdateWithVersionSql(ctx context.Context, tableFullName string, lockId string, exceptedVersion, newVersion storage.Version, lockInformation *storage.LockInformation) (string, []any) {
	insertSql := fmt.Sprintf("UPDATE %s SET version = $1, lock_information_json_string = $2 WHERE lock_id = $3 AND version = $5", tableFullName)
	return insertSql, []any{newVersion, lockInformation.ToJsonString(), lockId, lockInformation.OwnerId, exceptedVersion}
}

func (x *PostgresqlSqlProvider) CreateTableSql(ctx context.Context, tableFullName string) (string, []any) {
	createTableSql := `CREATE TABLE IF NOT EXISTS %s (
    lock_id VARCHAR(255) NOT NULL PRIMARY KEY,
    owner_id VARCHAR(255) NOT NULL,
    version BIGINT NOT NULL,
    lock_information_json_string VARCHAR(255) NOT NULL
)`
	return fmt.Sprintf(createTableSql, tableFullName), nil
}

func (x *PostgresqlSqlProvider) FindLockInformationJsonStringByIdSql(ctx context.Context, tableFullName string, lockId string) (string, []any) {
	getLockSql := fmt.Sprintf("SELECT lock_information_json_string FROM %s WHERE lock_id = $1", tableFullName)
	return getLockSql, []any{lockId}
}

func (x *PostgresqlSqlProvider) ListLockInformationJsonStringSql(ctx context.Context, tableFullName string) (string, []any) {
	sql := fmt.Sprintf("SELECT lock_information_json_string FROM %s", tableFullName)
	return sql, nil
}
