package main

import (
	"context"
	"fmt"

	postgresql_storage "github.com/storage-lock/go-postgresql-storage"
)

func main() {

	// 使用一个DSN形式的数据库连接字符串创建ConnectionManager
	host := "127.0.0.1"
	port := uint(5432)
	username := "postgres"
	passwd := "UeGqAm8CxYGldMDLoNNt"
	database := "postgres"
	connectionManager := postgresql_storage.NewPostgresqlConnectionManager(host, port, username, passwd, database)

	// 然后从这个ConnectionManager创建PostgreSQL Storage
	options := postgresql_storage.NewPostgresqlStorageOptions().SetConnectionManager(connectionManager)
	storage, err := postgresql_storage.NewPostgresqlStorage(context.Background(), options)
	if err != nil {
		panic(err)
	}
	fmt.Println(storage.GetName())

}
