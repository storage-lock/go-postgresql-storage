# Postgresql Storage 

# 一、这是什么
以Postgresql为存储引擎的[Storage](https://github.com/storage-lock/go-storage)实现，当前仓库为比较底层的存储层实现，你可以与[storage-lock](https://github.com/storage-lock/go-storage-lock)结合使用，或者这个项目[PostgreSQL-locks](https://github.com/storage-lock/go-postgresql-storage)里专门封装提供了一些PostgreSQL锁相关的更易用友好的API。


# 二、安装依赖

```bash
go get -u github.com/storage-lock/go-postgresql-storage
```

# 三、API Examples

## 3.1 从DSN创建PostgresqlStorage

在Golang的世界中连接数据库最常见的就是DSN，下面的例子演示了如何从一个DSN创建PostgresqlStorage： 

```go
package main

import (
	"context"
	"fmt"

	postgresql_storage "github.com/storage-lock/go-postgresql-storage"
)

func main() {

	// 使用一个DSN形式的数据库连接字符串创建ConnectionManager
	testDsn := "host=127.0.0.1 user=postgres password=UeGqAm8CxYGldMDLoNNt port=5432 dbname=postgres sslmode=disable"
	connectionManager := postgresql_storage.NewPostgresqlConnectionGetterFromDSN(testDsn)

	// 然后从这个ConnectionManager创建PostgreSQL Storage
	options := postgresql_storage.NewPostgresqlStorageOptions().SetConnectionManager(connectionManager)
	storage, err := postgresql_storage.NewPostgresqlStorage(context.Background(), options)
	if err != nil {
		panic(err)
	}
	fmt.Println(storage.GetName())

}
```


## 3.2 从连接属性（ip、端口、用户名、密码等等）中创建PostgresqlStorage

或者你的配置文件中存放的并不是DSN，而是零散的几个连接属性，下面是一个创建PostgresqlStorage的例子：

```go
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
```


## 3.3 从sql.DB创建PostgresqlStorage

或者现在你已经有从其它渠道创建的能够连接到Postgresql的sql.DB，则也可以从这个*sql.DB创建PostgresqlStorage

```go
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
```