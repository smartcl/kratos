#!/bin/bash

docker-compose -f quickstart.yml -f quickstart-standalone.yml up --force-recreate -d

#./kratos migrate sql "mysql://root:123456@tcp(127.0.0.1:13306)/kratos" -e -s
#
#
#mysql -uroot -p123456 -h127.0.0.1 -P3306 -e "DROP DATABASE IF EXISTS kratos;"
#mysql -uroot -p123456 -h127.0.0.1 -P3306 -e "CREATE DATABASE kratos;"
#./kratos -c config.yaml migrate sql -e --yes
#
#./kratos serve -c ./config.yaml
