#!/bin/bash

runSQL() {
  # docker exec -i "${CONTAINER_NAME}" \
  #             mysql \
  #             mysql_config_editor \
  #             set \
  #             --login-path=local \
  #             --host=localhost \
  #             --user="${MYSQL_USERNAME}" \
  #             --password="${MYSQL_PASSWORD}"

  local SQL=$1
  echo "exec sql - ${SQL}"

  start_time=$(date +%s)
  docker exec -i "${CONTAINER_NAME}" \
              mysql \
              -u"${MYSQL_USERNAME}" \
              -p"${MYSQL_PASSWORD}" \
              -e "${SQL}"
              # --login-path=local \

  end_time=$(date +%s)
  execution_time=$((end_time - start_time))
  echo "elapsed: ${execution_time} second"
}