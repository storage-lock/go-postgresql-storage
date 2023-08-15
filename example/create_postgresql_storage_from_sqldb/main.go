package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/storage-lock/go-storage"

	postgresql_storage "github.com/storage-lock/go-postgresql-storage"
)

func main() {

	// 使用一个DSN形式的数据库连接字符串创建ConnectionManager
	testDsn := "host=127.0.0.1 user=postgres password=UeGqAm8CxYGldMDLoNNt port=5432 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", testDsn)
	if err != nil {
		panic(err)
	}
	connectionManager := storage.NewFixedSqlDBConnectionManager(db)

	// 然后从这个ConnectionManager创建Postgresql Storage
	options := postgresql_storage.NewPostgresqlStorageOptions().SetConnectionManager(connectionManager)
	storage, err := postgresql_storage.NewPostgresqlStorage(context.Background(), options)
	if err != nil {
		panic(err)
	}
	fmt.Println(storage.GetName())

}
