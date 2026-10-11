package main

import (
	"context"
	"fmt"

	postgresql_storage "github.com/storage-lock/go-postgresql-storage"
)

func main() {
	// 使用一个DSN形式的数据库连接字符串创建ConnectionManager
	testDsn := "host=127.0.0.1 user=postgres password=UeGqAm8CxYGldMDLoNNt port=5432 dbname=postgres sslmode=disable"
	storage, err := newPostgresqlStorageFromDSN(context.Background(), testDsn)
	if err != nil {
		panic(err)
	}
	fmt.Println(storage.GetName())
}

func newPostgresqlStorageFromDSN(ctx context.Context, dsn string) (*postgresql_storage.PostgresqlStorage, error) {
	connectionManager := postgresql_storage.NewPostgresqlConnectionGetterFromDsn(dsn)
	options := postgresql_storage.NewPostgresqlStorageOptions().SetConnectionManager(connectionManager)
	return postgresql_storage.NewPostgresqlStorage(ctx, options)
}
