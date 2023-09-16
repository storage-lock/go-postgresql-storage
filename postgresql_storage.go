package postgresql_storage

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
	"github.com/storage-lock/go-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"strings"
	"time"
)

// PostgresqlStorage 基于Postgresql作为存储引擎
type PostgresqlStorage struct {

	// postgresql是sql storage的一种具体实现
	*sql_based_storage.SqlBasedStorage

	// 创建时所需的各种参数
	options *PostgresqlStorageOptions
}

var _ storage.Storage = &PostgresqlStorage{}

func NewPostgresqlStorage(ctx context.Context, options *PostgresqlStorageOptions) (*PostgresqlStorage, error) {

	// 参数校验
	if err := options.Check(); err != nil {
		return nil, err
	}

	// 先创建SQL Storage
	baseStorageOptions := sql_based_storage.NewSqlBasedStorageOptions().
		SetTableFullName(options.GetTableFullName()).
		SetConnectionManager(options.ConnectionManager).
		SetSqlProvider(NewPostgresqlSqlProvider())
	basedStorage, err := sql_based_storage.NewSqlBasedStorage(baseStorageOptions)
	if err != nil {
		return nil, err
	}

	// 然后再创建自己
	postgresqlStorage := &PostgresqlStorage{
		SqlBasedStorage: basedStorage,
		options:         options,
	}

	// 初始化
	err = postgresqlStorage.Init(ctx)
	if err != nil {
		return nil, err
	}

	return postgresqlStorage, nil
}

const PostgresqlStorageName = "postgresql-storage"

func (x *PostgresqlStorage) GetName() string {
	return PostgresqlStorageName
}

func (x *PostgresqlStorage) Init(ctx context.Context) (returnError error) {
	db, err := x.options.ConnectionManager.Take(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err := x.options.ConnectionManager.Return(ctx, db)
		if returnError == nil {
			returnError = err
		}
	}()

	// 这个逻辑跟通用的sql storage流程不太一样，所以这里就覆写这个Init方法
	// 如果设置了数据库的话需要切换数据库
	if x.options.Schema != "" {
		// 切换到schema，如果需要的话，但是schema不会自动创建，需要使用者自己创建，会自动创建的只有存放锁信息的表
		_, err = db.ExecContext(ctx, fmt.Sprintf("SET search_path TO %s ", x.options.Schema))
		if err != nil {
			return err
		}
	}

	// 创建存储锁信息需要的表
	createTableSql := `CREATE TABLE IF NOT EXISTS %s (
    lock_id VARCHAR(255) NOT NULL PRIMARY KEY,
    owner_id VARCHAR(255) NOT NULL, 
    version BIGINT NOT NULL,
    lock_information_json_string VARCHAR(255) NOT NULL
)`
	_, err = db.ExecContext(ctx, fmt.Sprintf(createTableSql, x.options.GetTableName()))
	if err != nil {
		return err
	}

	return nil
}

// duplicate key value violates unique constraint

func (x *PostgresqlStorage) CreateWithVersion(ctx context.Context, lockId string, version storage.Version, lockInformation *storage.LockInformation) (returnError error) {

	returnError = x.SqlBasedStorage.CreateWithVersion(ctx, lockId, version, lockInformation)
	if returnError != nil {
		// 把重复转为版本miss
		if strings.Contains(returnError.Error(), "duplicate key value violates unique constraint") {
			return storage_lock.ErrVersionMiss
		}
	}
	return returnError
}

func (x *PostgresqlStorage) GetTime(ctx context.Context) (time.Time, error) {

	db, err := x.options.ConnectionManager.Take(ctx)
	if err != nil {
		return time.Time{}, err
	}
	defer func() {
		_ = x.options.ConnectionManager.Return(ctx, db)
	}()

	var zero time.Time
	rs, err := db.Query("SELECT CURRENT_TIMESTAMP")
	if err != nil {
		return zero, err
	}
	defer func() {
		_ = rs.Close()
	}()
	if !rs.Next() {
		return zero, ErrQueryPostgresqlServerTime
	}
	var databaseTime time.Time
	err = rs.Scan(&databaseTime)
	if err != nil {
		return zero, err
	}

	return databaseTime, nil
}
