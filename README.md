# Postgresql Storage 


```bash
go get -u github.com/storage-lock/go-postgresql-storage
```



## 3.2 Postgresql

### 3.2.1 快速开始

```go
package main

import (
	"context"
	"fmt"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"strings"
	"sync"
	"time"
)

func main() {

	// Docker启动Postgresql：
	// docker run -d --name storage-lock-postgres -p 5432:5432 -e POSTGRES_PASSWORD=UeGqAm8CxYGldMDLoNNt postgres:14

	// DSN的写法参考驱动的支持：https://github.com/lib/pq
	dsn := "host=192.168.128.206 user=postgres password=UeGqAm8CxYGldMDLoNNt port=5432 dbname=postgres sslmode=disable"

	// 这个是最为重要的，通常是要锁住的资源的名称
	lockId := "must-serial-operation-resource-foo"

	// 第一步创建一把分布式锁
	lock, err := storage_lock.NewPostgreSQLStorageLock(context.Background(), lockId, dsn)
	if err != nil {
		fmt.Printf("[ %s ] Create Lock Failed: %v\n", time.Now().Format("2006-01-02 15:04:05"), err)
		return
	}

	// 第二步使用这把锁，这里就模拟多个节点竞争执行的情况，他们会线程安全的往resource里写数据
	resource := strings.Builder{}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		workerId := fmt.Sprintf("worker-%d", i)
		wg.Add(1)
		go func() {
			defer wg.Done()

			// 获取锁
			err := lock.Lock(context.Background(), workerId)
			if err != nil {
				fmt.Printf("[ %s ] workerId = %s, lock failed: %v \n", time.Now().Format("2006-01-02 15:04:05"), workerId, err)
				return
			}
			// 退出的时候释放锁
			defer func() {
				err := lock.UnLock(context.Background(), workerId)
				if err != nil {
					fmt.Printf("[ %s ] workerId = %s, unlock failed: %v \n", time.Now().Format("2006-01-02 15:04:05"), workerId, err)
					return
				}
			}()

			// 假装有耗时的操作
			fmt.Printf("[ %s ] workerId = %s, begin write resource \n", time.Now().Format("2006-01-02 15:04:05"), workerId)
			time.Sleep(time.Second * 3)
			// 接下来是操作竞态资源
			resource.WriteString(workerId)
			fmt.Printf("[ %s ] workerId = %s, write resource done \n", time.Now().Format("2006-01-02 15:04:05"), workerId)
			resource.WriteString("\n")

		}()
	}
	wg.Wait()

	// 观察最终的输出是否和日志一致
	fmt.Printf("[ %s ] Resource: \n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(resource.String())

	// Output:
	// [ 2023-03-13 00:29:37 ] workerId = worker-3, begin write resource
	// [ 2023-03-13 00:29:40 ] workerId = worker-3, write resource done
	// [ 2023-03-13 00:29:40 ] workerId = worker-5, begin write resource
	// [ 2023-03-13 00:29:43 ] workerId = worker-5, write resource done
	// [ 2023-03-13 00:29:43 ] workerId = worker-8, begin write resource
	// [ 2023-03-13 00:29:46 ] workerId = worker-8, write resource done
	// [ 2023-03-13 00:29:46 ] workerId = worker-6, begin write resource
	// [ 2023-03-13 00:29:49 ] workerId = worker-6, write resource done
	// [ 2023-03-13 00:29:50 ] workerId = worker-2, begin write resource
	// [ 2023-03-13 00:29:53 ] workerId = worker-2, write resource done
	// [ 2023-03-13 00:29:56 ] workerId = worker-0, begin write resource
	// [ 2023-03-13 00:29:59 ] workerId = worker-0, write resource done
	// [ 2023-03-13 00:30:00 ] workerId = worker-1, begin write resource
	// [ 2023-03-13 00:30:03 ] workerId = worker-1, write resource done
	// [ 2023-03-13 00:30:04 ] workerId = worker-4, begin write resource
	// [ 2023-03-13 00:30:07 ] workerId = worker-4, write resource done
	// [ 2023-03-13 00:30:08 ] workerId = worker-9, begin write resource
	// [ 2023-03-13 00:30:11 ] workerId = worker-9, write resource done
	// [ 2023-03-13 00:30:14 ] workerId = worker-7, begin write resource
	// [ 2023-03-13 00:30:18 ] workerId = worker-7, write resource done
	// [ 2023-03-13 00:30:18 ] Resource:
	// worker-3
	// worker-5
	// worker-8
	// worker-6
	// worker-2
	// worker-0
	// worker-1
	// worker-4
	// worker-9
	// worker-7

}

```

### 3.2.2 详细配置

```go
package main

import (
	"context"
	"fmt"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"strings"
	"sync"
	"time"
)

func main() {

	// Docker启动Postgresql：
	// docker run -d --name storage-lock-postgres -p 5432:5432 -e POSTGRES_PASSWORD=UeGqAm8CxYGldMDLoNNt postgres:14

	// DSN的写法参考驱动的支持：https://github.com/lib/pq
	dsn := "host=192.168.128.206 user=postgres password=UeGqAm8CxYGldMDLoNNt port=5432 dbname=postgres sslmode=disable"

	// 第一步先配置存储介质相关的参数，包括如何连接到这个数据库，连接上去之后锁的信息存储到哪里等等
	// 配置如何连接到数据库
	connectionGetter := storage_lock.NewPostgreSQLStorageConnectionGetterFromDSN(dsn)
	storageOptions := &storage_lock.PostgreSQLStorageOptions{
		// 数据库连接获取方式，可以使用内置的从DSN获取连接，也可以自己实现接口决定如何连接
		ConnectionGetter: connectionGetter,
		// 选择锁信息存放在哪个schema下，默认为public
		Schema: "public",
		// 锁的信息是存储在哪张表中的，不设置的话默认为storage_lock
		TableName: "storage_lock_table",
	}
	storage, err := storage_lock.NewPostgreSQLStorage(context.Background(), storageOptions)
	if err != nil {
		fmt.Println("Create Storage Failed： " + err.Error())
		return
	}

	// 第二步配置锁的参数，在上面创建的Storage的上创建一把锁
	lockOptions := &storage_lock.StorageLockOptions{
		// 这个是最为重要的，通常是要锁住的资源的名称
		LockId:                "must-serial-operation-resource-foo",
		LeaseExpireAfter:      time.Second * 30,
		LeaseRefreshInterval:  time.Second * 5,
		VersionMissRetryTimes: 3,
	}
	lock := storage_lock.NewStorageLock(storage, lockOptions)

	// 第三步开始使用锁，模拟多个节点竞争同一个锁使用的情况
	resource := strings.Builder{}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		workerId := fmt.Sprintf("worker-%d", i)
		wg.Add(1)
		go func() {
			defer wg.Done()

			// 获取锁
			err := lock.Lock(context.Background(), workerId)
			if err != nil {
				fmt.Printf("[ %s ] workerId = %s, lock failed: %v \n", time.Now().Format("2006-01-02 15:04:05"), workerId, err)
				return
			}
			// 退出的时候释放锁
			defer func() {
				err := lock.UnLock(context.Background(), workerId)
				if err != nil {
					fmt.Printf("[ %s ] workerId = %s, unlock failed: %v \n", time.Now().Format("2006-01-02 15:04:05"), workerId, err)
					return
				}
			}()

			// 假装有耗时的操作
			fmt.Printf("[ %s ] workerId = %s, begin write resource \n", time.Now().Format("2006-01-02 15:04:05"), workerId)
			time.Sleep(time.Second * 3)
			// 接下来是操作竞态资源
			resource.WriteString(workerId)
			fmt.Printf("[ %s ] workerId = %s, write resource done \n", time.Now().Format("2006-01-02 15:04:05"), workerId)
			resource.WriteString("\n")

		}()
	}
	wg.Wait()

	// 观察最终的输出是否和日志一致
	fmt.Printf("[ %s ] Resource: \n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(resource.String())

	// Output:
	// [ 2023-03-13 00:33:38 ] workerId = worker-0, begin write resource
	// [ 2023-03-13 00:33:41 ] workerId = worker-0, write resource done
	// [ 2023-03-13 00:33:42 ] workerId = worker-3, begin write resource
	// [ 2023-03-13 00:33:45 ] workerId = worker-3, write resource done
	// [ 2023-03-13 00:33:45 ] workerId = worker-6, begin write resource
	// [ 2023-03-13 00:33:48 ] workerId = worker-6, write resource done
	// [ 2023-03-13 00:33:49 ] workerId = worker-5, begin write resource
	// [ 2023-03-13 00:33:52 ] workerId = worker-5, write resource done
	// [ 2023-03-13 00:33:53 ] workerId = worker-2, begin write resource
	// [ 2023-03-13 00:33:56 ] workerId = worker-2, write resource done
	// [ 2023-03-13 00:33:57 ] workerId = worker-8, begin write resource
	// [ 2023-03-13 00:34:00 ] workerId = worker-8, write resource done
	// [ 2023-03-13 00:34:01 ] workerId = worker-4, begin write resource
	// [ 2023-03-13 00:34:04 ] workerId = worker-4, write resource done
	// [ 2023-03-13 00:34:04 ] workerId = worker-1, begin write resource
	// [ 2023-03-13 00:34:07 ] workerId = worker-1, write resource done
	// [ 2023-03-13 00:34:08 ] workerId = worker-9, begin write resource
	// [ 2023-03-13 00:34:11 ] workerId = worker-9, write resource done
	// [ 2023-03-13 00:34:11 ] workerId = worker-7, begin write resource
	// [ 2023-03-13 00:34:14 ] workerId = worker-7, write resource done
	// [ 2023-03-13 00:34:14 ] Resource:
	// worker-0
	// worker-3
	// worker-6
	// worker-5
	// worker-2
	// worker-8
	// worker-4
	// worker-1
	// worker-9
	// worker-7

}

```