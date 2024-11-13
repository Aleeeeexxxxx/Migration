#!/bin/bash

set -euo pipefail

declare TABLE_NAME='models'

declare CONTAINER_NAME='mysql'
declare MYSQL_ROOT_USERNAME='root'
declare MYSQL_ROOT_PASSWORD='root'

declare MYSQL_DEFAULT_USERNAME='origin'
declare MYSQL_DEFAULT_PASSWORD='origin'

# run mysql
if [ ! "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
    if [ "$(docker ps -aq -f status=exited -f name=$CONTAINER_NAME)" ]; then
        docker rm $CONTAINER_NAME
    fi
    make build_mysql
    docker run -d \
               -p 3306:3306 \
               --name "${CONTAINER_NAME}" \
               -e MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD}" \
               alex-mysql:latest
    sleep 10 # wait for mysql up
fi

echo 'container mysql started'

declare MYSQL_USERNAME="${MYSQL_ROOT_USERNAME}"
declare MYSQL_PASSWORD="${MYSQL_ROOT_USERNAME}"
source ./script/common.sh

# create users
runSQL "create user if not exists '${MYSQL_DEFAULT_USERNAME}'@'%' IDENTIFIED BY '${MYSQL_DEFAULT_PASSWORD}';"

# create databases and grant privileges
db_to_create=(
  "migration_origin"
  "migration_migrated1"
  "migration_migrated2"
  "migration_migrated3"
)

for item in "${db_to_create[@]}"; do
    runSQL "CREATE DATABASE IF NOT EXISTS ${item} ;"
    runSQL "CREATE TABLE IF NOT EXISTS ${item}.${TABLE_NAME} (
              id VARCHAR(191) NOT NULL PRIMARY KEY,
              msg VARCHAR(255) DEFAULT NULL,
              updated_at BIGINT DEFAULT NULL
            );"
    runSQL "GRANT ALL PRIVILEGES ON ${item}.* TO '${MYSQL_DEFAULT_USERNAME}'@'%';"
    runSQL "GRANT FILE ON *.* TO 'origin'@'%';"
done

# restrict root 
runSQL "DELETE FROM mysql.user WHERE User = ${MYSQL_ROOT_USERNAME} AND Host != 'localhost';"
runSQL "UPDATE mysql.user SET Host = 'localhost' WHERE User = ${MYSQL_ROOT_USERNAME};"

# flush
runSQL "FLUSH PRIVILEGES;"

echo 'mysql has been deployed successfully'






