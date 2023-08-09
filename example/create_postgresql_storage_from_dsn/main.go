package main

import (
	"context"
	"fmt"

	postgresql_storage "github.com/storage-lock/go-postgresql-storage"
)

func main() {

	// 使用一个DSN形式的数据库连接字符串创建ConnectionManager
	testDsn := "host=127.0.0.1 user=postgres password=UeGqAm8CxYGldMDLoNNt port=5432 dbname=postgres sslmode=disable"
	connectionManager := postgresql_storage.NewPostgreSQLConnectionGetterFromDSN(testDsn)

	// 然后从这个ConnectionManager创建PostgreSQL Storage
	options := postgresql_storage.NewPostgreSQLStorageOptions().SetConnectionManager(connectionManager)
	storage, err := postgresql_storage.NewPostgreSQLStorage(context.Background(), options)
	if err != nil {
		panic(err)
	}
	fmt.Println(storage.GetName())

}
