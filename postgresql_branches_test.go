package postgresql_storage

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/storage-lock/go-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"github.com/stretchr/testify/assert"
)

type takeErrorManager struct{}

func (takeErrorManager) Name() string { return "take-error" }
func (takeErrorManager) Take(context.Context) (*sql.DB, error) {
	return nil, errors.New("take failed")
}
func (takeErrorManager) Return(context.Context, *sql.DB) error { return nil }
func (takeErrorManager) Shutdown(context.Context) error        { return nil }

type returnErrorManager struct{ db *sql.DB }

func (returnErrorManager) Name() string { return "return-error" }
func (x returnErrorManager) Take(context.Context) (*sql.DB, error) {
	return x.db, nil
}
func (returnErrorManager) Return(context.Context, *sql.DB) error {
	return errors.New("return failed")
}
func (returnErrorManager) Shutdown(context.Context) error { return nil }

type defaultDriver struct{}

func (defaultDriver) Open(string) (driver.Conn, error) { return defaultConn{}, nil }

type defaultConn struct{}

func (defaultConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("unsupported") }
func (defaultConn) Close() error                        { return nil }
func (defaultConn) Begin() (driver.Tx, error)           { return nil, errors.New("unsupported") }

func init() {
	sql.Register("postgres-coverage", defaultDriver{})
}

func TestPostgresqlConnectionManager(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectClose()

	fromDSN := NewPostgresqlConnectionGetterFromDsn("host=h user=u")
	assert.Equal(t, PostgresqlConnectionManagerName, fromDSN.Name())
	assert.Equal(t, "host=h user=u", fromDSN.GetDSN())
	fromDSN.open = func(string, string) (*sql.DB, error) { return nil, errors.New("open failed") }
	_, err = fromDSN.Take(context.Background())
	assert.EqualError(t, err, "open failed")
	_, err = fromDSN.Take(context.Background())
	assert.EqualError(t, err, "open failed")

	preset := &PostgresqlConnectionManager{err: errors.New("preset")}
	_, err = preset.Take(context.Background())
	assert.EqualError(t, err, "preset")

	fromDB := NewPostgresqlConnectionGetterFromSqlDb(db)
	got, err := fromDB.Take(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, db, got)
	assert.NoError(t, fromDB.Return(context.Background(), got))
	assert.NoError(t, fromDB.Shutdown(context.Background()))

	fromConfig := NewPostgresqlConnectionManager("h", 5432, "u", "p", "db")
	assert.Contains(t, fromConfig.GetDSN(), "host=h")
	fromConfig.open = func(string, string) (*sql.DB, error) { return db, nil }
	got, err = fromConfig.Take(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, db, got)
	assert.NoError(t, (&PostgresqlConnectionManager{}).Shutdown(context.Background()))

	defaultManager := NewPostgresqlConnectionGetterFromDsn("host=127.0.0.1 port=1 user=u dbname=db sslmode=disable")
	opened, err := defaultManager.Take(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, opened)
	assert.NoError(t, defaultManager.Shutdown(context.Background()))

	hooked := NewPostgresqlConnectionGetterFromDsn("host=127.0.0.1")
	hooked.open = func(driverName, dsn string) (*sql.DB, error) {
		return sql.Open("postgres-coverage", dsn)
	}
	opened, err = hooked.Take(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, opened)
	assert.NoError(t, hooked.Shutdown(context.Background()))
}

func TestPostgresqlSqlProvider(t *testing.T) {
	provider := NewPostgresqlSqlProvider()
	ctx := context.Background()
	info := &storage.LockInformation{OwnerId: "owner"}

	sqlText, args := provider.NowTimestampSql(ctx, "public.locks")
	assert.Equal(t, "SELECT CURRENT_TIMESTAMP", sqlText)
	assert.Nil(t, args)

	sqlText, args = provider.SelectLockInformationJsonStringSql(ctx, "public.locks", "lock")
	assert.Contains(t, sqlText, "$1")
	assert.Equal(t, []any{"lock"}, args)

	sqlText, args = provider.DeleteWithVersionSql(ctx, "public.locks", "lock", 3, info)
	assert.Contains(t, sqlText, "DELETE")
	assert.Equal(t, []any{"lock", "owner", storage.Version(3)}, args)

	sqlText, args = provider.CreateWithVersionSql(ctx, "public.locks", "lock", 1, info)
	assert.Contains(t, sqlText, "INSERT")
	assert.Len(t, args, 4)

	sqlText, args = provider.UpdateWithVersionSql(ctx, "public.locks", "lock", 1, 2, info)
	assert.Contains(t, sqlText, "UPDATE")
	assert.Equal(t, storage.Version(2), args[0])

	sqlText, args = provider.CreateTableSql(ctx, "public.locks")
	assert.Contains(t, sqlText, "CREATE TABLE")
	assert.Nil(t, args)

	sqlText, args = provider.FindLockInformationJsonStringByIdSql(ctx, "public.locks", "lock")
	assert.Contains(t, sqlText, "lock_id = $1")
	assert.Equal(t, []any{"lock"}, args)

	sqlText, args = provider.ListLockInformationJsonStringSql(ctx, "public.locks")
	assert.Contains(t, sqlText, "SELECT lock_information_json_string")
	assert.Nil(t, args)
}

func TestNewPostgresqlStorage_Branches(t *testing.T) {
	_, err := NewPostgresqlStorage(context.Background(), NewPostgresqlStorageOptions())
	assert.EqualError(t, err, "ConnectionManager can not nil")

	opts := NewPostgresqlStorageOptions()
	opts.ConnectionManager = storage.NewFixedSqlDBConnectionManager(&sql.DB{})
	opts.Schema = ""
	opts.TableName = ""
	assert.NoError(t, opts.Check())
	assert.Equal(t, DefaultPostgresqlStorageSchema, opts.Schema)
	assert.Equal(t, storage.DefaultStorageTableName, opts.TableName)

	_, err = NewPostgresqlStorage(context.Background(), &PostgresqlStorageOptions{
		ConnectionManager: takeErrorManager{},
		Schema:            "public",
		TableName:         "locks",
	})
	assert.EqualError(t, err, "take failed")

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectExec("SET search_path").WillReturnError(errors.New("schema failed"))
	_, err = NewPostgresqlStorage(context.Background(), &PostgresqlStorageOptions{
		ConnectionManager: storage.NewFixedSqlDBConnectionManager(db),
		Schema:            "custom",
		TableName:         "locks",
	})
	assert.EqualError(t, err, "schema failed")

	db2, mock2, err := sqlmock.New()
	assert.NoError(t, err)
	defer db2.Close()
	mock2.ExpectExec("SET search_path").WillReturnResult(sqlmock.NewResult(0, 0))
	mock2.ExpectExec("CREATE TABLE").WillReturnError(errors.New("create failed"))
	mock2.ExpectClose()
	_, err = NewPostgresqlStorage(context.Background(), &PostgresqlStorageOptions{
		ConnectionManager: returnErrorManager{db: db2},
		Schema:            "public",
		TableName:         "locks",
	})
	assert.Error(t, err)

	db3, mock3, err := sqlmock.New()
	assert.NoError(t, err)
	defer db3.Close()
	mock3.ExpectExec("SET search_path").WillReturnResult(sqlmock.NewResult(0, 0))
	mock3.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))
	s, err := NewPostgresqlStorage(context.Background(), &PostgresqlStorageOptions{
		ConnectionManager: storage.NewFixedSqlDBConnectionManager(db3),
		Schema:            "public",
		TableName:         "locks",
	})
	assert.NoError(t, err)

	info := &storage.LockInformation{OwnerId: "owner"}
	mock3.ExpectExec("INSERT INTO").WillReturnError(errors.New("duplicate key value violates unique constraint"))
	assert.ErrorIs(t, s.CreateWithVersion(context.Background(), "lock", 1, info), storage_lock.ErrVersionMiss)
	mock3.ExpectExec("INSERT INTO").WillReturnError(errors.New("other"))
	assert.EqualError(t, s.CreateWithVersion(context.Background(), "lock", 1, info), "other")
	mock3.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	assert.NoError(t, s.CreateWithVersion(context.Background(), "lock", 1, info))

	_, err = s.GetTime(context.Background())
	assert.Error(t, err)
	mock3.ExpectQuery("SELECT CURRENT_TIMESTAMP").WillReturnError(errors.New("query failed"))
	_, err = s.GetTime(context.Background())
	assert.EqualError(t, err, "query failed")
	mock3.ExpectQuery("SELECT CURRENT_TIMESTAMP").WillReturnRows(sqlmock.NewRows([]string{"now"}))
	_, err = s.GetTime(context.Background())
	assert.ErrorIs(t, err, ErrQueryPostgresqlServerTime)
	now := time.Now()
	mock3.ExpectQuery("SELECT CURRENT_TIMESTAMP").WillReturnRows(sqlmock.NewRows([]string{"now"}).AddRow(now))
	got, err := s.GetTime(context.Background())
	assert.NoError(t, err)
	assert.False(t, got.IsZero())
	mock3.ExpectQuery("SELECT CURRENT_TIMESTAMP").WillReturnRows(sqlmock.NewRows([]string{"now"}).AddRow("bad"))
	_, err = s.GetTime(context.Background())
	assert.Error(t, err)

	takeFailed := &PostgresqlStorage{options: &PostgresqlStorageOptions{ConnectionManager: takeErrorManager{}}}
	_, err = takeFailed.GetTime(context.Background())
	assert.EqualError(t, err, "take failed")

	_, err = (&PostgresqlStorage{}).Watch(context.Background(), "")
	assert.Error(t, err)
}
