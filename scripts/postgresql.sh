#!/usr/bin/env bash

# 删除可能存在的容器
docker rm -f storage-lock-postgres

# 启动测试使用的Postgresql数据库
docker run -d --name storage-lock-postgres -p 5432:5432 -e POSTGRES_PASSWORD=UeGqAm8CxYGldMDLoNNt postgres:14

# 设置环境变量
export STORAGE_LOCK_POSTGRESQL_DSN="host=127.0.0.1 user=postgres password=UeGqAm8CxYGldMDLoNNt port=5432 dbname=postgres sslmode=disable"