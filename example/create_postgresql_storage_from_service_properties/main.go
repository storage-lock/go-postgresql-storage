package main

import (
	"context"
	"fmt"

	postgresql_storage "github.com/storage-lock/go-postgresql-storage"
)

func main() {
	storage, err := newPostgresqlStorageFromProperties(context.Background(), "127.0.0.1", 5432, "postgres", "UeGqAm8CxYGldMDLoNNt", "postgres")
	if err != nil {
		panic(err)
	}
	fmt.Println(storage.GetName())
}

func newPostgresqlStorageFromProperties(ctx context.Context, host string, port uint, username, passwd, database string) (*postgresql_storage.PostgresqlStorage, error) {
	connectionManager := postgresql_storage.NewPostgresqlConnectionManager(host, port, username, passwd, database)
	options := postgresql_storage.NewPostgresqlStorageOptions().SetConnectionManager(connectionManager)
	return postgresql_storage.NewPostgresqlStorage(ctx, options)
}
