#!/bin/bash

GOOS=linux GOARCH=arm64 go build
docker cp kratos ory-kratos-hydra-integration-demo-kratos-1:/usr/bin/kratos
docker commit -a "cl" -m "保存当前状态" ory-kratos-hydra-integration-demo-kratos-1 oryd/kratos:v1.2.0.1
