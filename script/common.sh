#!/bin/bash

runSQL() {
  set +x

  local SQL=$1
  echo "exec sql - ${SQL}"

  start_time=$(date +%s)
  docker exec -i "${CONTAINER_NAME}" \
              mysql \
              -u"${MYSQL_USERNAME}" \
              -p"${MYSQL_PASSWORD}" \
              -e "${SQL}"

  end_time=$(date +%s)
  execution_time=$((end_time - start_time))
  echo "elapsed: ${execution_time} second"
    
  set -x
}