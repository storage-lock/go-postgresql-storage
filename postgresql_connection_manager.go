package postgresql_storage

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/storage-lock/go-storage"
)

const DefaultPostgreSQLStorageSchema = "public"

type PostgreSQLConnectionManager struct {

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

var _ storage.ConnectionManager[*sql.DB] = &PostgreSQLConnectionManager{}

// NewPostgreSQLConnectionGetterFromDSN 从DSN创建PostgreSQL连接
func NewPostgreSQLConnectionGetterFromDSN(dsn string) *PostgreSQLConnectionManager {
	return &PostgreSQLConnectionManager{
		DSN: dsn,
	}
}

// NewPostgreSQLConnectionGetter 从服务器属性创建数据库连接
func NewPostgreSQLConnectionGetter(host string, port uint, user, passwd, databaseName string) *PostgreSQLConnectionManager {
	return &PostgreSQLConnectionManager{
		Host:         host,
		Port:         port,
		User:         user,
		Passwd:       passwd,
		DatabaseName: databaseName,
	}
}

func (x *PostgreSQLConnectionManager) Name() string {
	return "postgresql-connection-manager"
}

func (x *PostgreSQLConnectionManager) GetDSN() string {
	if x.DSN != "" {
		return x.DSN
	}
	return fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable", x.Host, x.User, x.Passwd, x.Port, x.DatabaseName)
}

// Take 获取到数据库的连接
func (x *PostgreSQLConnectionManager) Take(ctx context.Context) (*sql.DB, error) {
	x.once.Do(func() {
		db, err := sql.Open("postgres", x.GetDSN())
		if err != nil {
			x.err = err
			return
		}
		x.db = db
	})
	return x.db, x.err
}

func (x *PostgreSQLConnectionManager) Return(ctx context.Context, connection *sql.DB) error {
	return nil
}

func (x *PostgreSQLConnectionManager) Shutdown(ctx context.Context) error {
	if x.db != nil {
		return x.db.Close()
	}
	return nil
}
