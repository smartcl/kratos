#!/bin/bash
docker rm -f kratos-kratos-1 kratos-kratos-migrate-1 kratos-kratos-selfservice-ui-node-1 kratos-mailslurper-1 kratos-mysqld-1
docker volume prune
