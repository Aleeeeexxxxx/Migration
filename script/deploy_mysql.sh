#!/bin/bash

set -euo pipefail

declare CONTAINER_NAME='mysql'
declare MYSQL_ROOT_USERNAME='root'
declare MYSQL_ROOT_PASSWORD='root'

declare MYSQL_USERNAME='origin'
declare MYSQL_PASSWORD='origin'

# run mysql
if [ ! "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
    if [ "$(docker ps -aq -f status=exited -f name=$CONTAINER_NAME)" ]; then
        docker rm $CONTAINER_NAME
    fi
    docker run -p 3306:3306 \
               --name "${CONTAINER_NAME}" \
               -e MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD}" \
               mysql
fi

echo 'container mysql started'

runSQL() {
  local SQL=$1
  echo "exec sql - ${SQL}"
  docker exec -i "${CONTAINER_NAME}" \
              mysql \
              -u"${MYSQL_ROOT_USERNAME}" \
              -p"${MYSQL_ROOT_USERNAME}" \
              -e "${SQL}"
}

# restrict root privilege
runSQL "UPDATE mysql.user SET Host='localhost' WHERE User='root';"

# create users
runSQL "create user if not exists '${MYSQL_USERNAME}'@'%' IDENTIFIED BY '${MYSQL_PASSWORD}';"

# create databases and grant privileges
db_to_create=(
  "migration_origin"
  "migration_migrated1"
  "migration_migrated2"
  "migration_migrated3"
)

for item in "${db_to_create[@]}"; do
    runSQL "create database if not exists ${item} ;"
    runSQL "GRANT ALL PRIVILEGES ON ${item}.* TO '${MYSQL_USERNAME}'@'%';"
done

# flush
runSQL "FLUSH PRIVILEGES;"

echo 'mysql has been deployed successfully'






