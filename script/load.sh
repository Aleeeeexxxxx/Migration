#!/bin/bash

set -euox pipefail

declare USE_LOADER=false

declare GENERATED_NUM=10000
declare OUTPUT_FILE='data.csv'
declare TABLE_NAME='migration_origin.models'

declare CONTAINER_NAME='mysql'
declare MYSQL_USERNAME='origin'
declare MYSQL_PASSWORD='origin'

rm -f ${OUTPUT_FILE}

# generate data
docker exec -i "${CONTAINER_NAME}" \
              generator \
              -o "${OUTPUT_FILE}" \
              -n "${GENERATED_NUM}"

if [ "$USE_LOADER" = true ]; then
    echo "using loader"
    docker exec -i "${CONTAINER_NAME}" \
                  loader \
                  -f "${OUTPUT_FILE}" 
else
    echo "using load data infile"

    SQL="LOAD DATA INFILE '${OUTPUT_FILE}'
         INTO TABLE ${TABLE_NAME}
         FIELDS TERMINATED BY ','
         LINES TERMINATED BY '\n'
         (ID, MSG, UPDATED_AT);"
    start_time=$(date +%s)

    docker exec -i "${CONTAINER_NAME}" \
                mysql \
                -u"${MYSQL_USERNAME}" \
                -p"${MYSQL_PASSWORD}" \
                -e "${SQL}"

    end_time=$(date +%s)
    execution_time=$((end_time - start_time))
    echo "elapsed：${execution_time} second"
fi

