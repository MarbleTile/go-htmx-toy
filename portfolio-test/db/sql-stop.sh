#!/bin/bash

docker kill go-test-mysql
docker container prune -f
