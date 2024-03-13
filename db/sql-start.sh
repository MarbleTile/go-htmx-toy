#!/bin/bash

docker run --network mysql-net --name go-test-mysql -e MYSQL_ROOT_PASSWORD=cumdump -d mysql:latest
