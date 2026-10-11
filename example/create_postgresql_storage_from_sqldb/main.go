package main

import (
	"context"
	"database/sql"
	"fmt"

	postgresql_storage "github.com/storage-lock/go-postgresql-storage"
	"github.com/storage-lock/go-storage"
)

var sqlOpen = sql.Open

func main() {
	// 使用一个DSN形式的数据库连接字符串创建ConnectionManager
	testDsn := "host=127.0.0.1 user=postgres password=UeGqAm8CxYGldMDLoNNt port=5432 dbname=postgres sslmode=disable"
	db, err := sqlOpen("postgres", testDsn)
	if err != nil {
		panic(err)
	}
	storage, err := newPostgresqlStorageFromDB(context.Background(), db)
	if err != nil {
		panic(err)
	}
	fmt.Println(storage.GetName())
}

func newPostgresqlStorageFromDB(ctx context.Context, db *sql.DB) (*postgresql_storage.PostgresqlStorage, error) {
	connectionManager := storage.NewFixedSqlDBConnectionManager(db)
	options := postgresql_storage.NewPostgresqlStorageOptions().SetConnectionManager(connectionManager)
	return postgresql_storage.NewPostgresqlStorage(ctx, options)
}
