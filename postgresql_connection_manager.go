package postgresql_storage

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/storage-lock/go-storage"
)

// DefaultPostgresqlStorageSchema 默认的schema
const DefaultPostgresqlStorageSchema = "public"

// PostgresqlConnectionManager Postgresql的连接管理器
type PostgresqlConnectionManager struct {

	// 主机的名字
	Host string

	// 主机的端口
	Port uint

	// 用户名
	User string

	// 密码
	Passwd string

	DatabaseName string

	// DSN
	// Example: "host=192.168.128.206 user=postgres password=123456 port=5432 dbname=postgres sslmode=disable"
	DSN string

	// 初始化好的数据库实例
	db   *sql.DB
	err  error
	once sync.Once
}

var _ storage.ConnectionManager[*sql.DB] = &PostgresqlConnectionManager{}

// NewPostgresqlConnectionGetterFromDSN 从DSN创建PostgreSQL连接
func NewPostgresqlConnectionGetterFromDSN(dsn string) *PostgresqlConnectionManager {
	return &PostgresqlConnectionManager{
		DSN: dsn,
	}
}

// NewPostgresqlConnectionManager 从服务器属性创建数据库连接
func NewPostgresqlConnectionManager(host string, port uint, user, passwd, databaseName string) *PostgresqlConnectionManager {
	return &PostgresqlConnectionManager{
		Host:         host,
		Port:         port,
		User:         user,
		Passwd:       passwd,
		DatabaseName: databaseName,
	}
}

// NewPostgresqlConnectionGetterFromSqlDb 从一个已经存在的*sql.DB创建连接管理器
func NewPostgresqlConnectionGetterFromSqlDb(db *sql.DB) *PostgresqlConnectionManager {
	return &PostgresqlConnectionManager{
		db: db,
	}
}

const PostgreSQLConnectionManagerName = "postgresql-connection-manager"

func (x *PostgresqlConnectionManager) Name() string {
	return PostgreSQLConnectionManagerName
}

func (x *PostgresqlConnectionManager) GetDSN() string {
	if x.DSN != "" {
		return x.DSN
	}
	return fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable", x.Host, x.User, x.Passwd, x.Port, x.DatabaseName)
}

// Take 获取到数据库的连接
func (x *PostgresqlConnectionManager) Take(ctx context.Context) (*sql.DB, error) {
	x.once.Do(func() {
		if x.err != nil {
			return
		}
		db, err := sql.Open("postgres", x.GetDSN())
		if err != nil {
			x.err = err
			return
		}
		x.db = db
	})
	return x.db, x.err
}

func (x *PostgresqlConnectionManager) Return(ctx context.Context, connection *sql.DB) error {
	return nil
}

func (x *PostgresqlConnectionManager) Shutdown(ctx context.Context) error {
	if x.db != nil {
		return x.db.Close()
	}
	return nil
}
