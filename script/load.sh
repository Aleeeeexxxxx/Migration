#!/bin/bash

set -euo pipefail

declare USE_LOADER=false

declare GENERATED_NUM=1000000
declare OUTPUT_FILE='/var/lib/mysql-files/data.csv'
declare TABLE_NAME='migration_origin.models'

declare CONTAINER_NAME='mysql'
declare MYSQL_USERNAME='origin'
declare MYSQL_PASSWORD='origin'

source ./script/common.sh

# remove old data file
rm -f ${OUTPUT_FILE}

# generate data
docker exec -i "${CONTAINER_NAME}" \
            generator \
            -o "${OUTPUT_FILE}" \
            -n "${GENERATED_NUM}"

echo "data generated. total=${GENERATED_NUM}"

# remove all data
runSQL "DELETE FROM ${TABLE_NAME} WHERE 1=1;"

# load
if [ "$USE_LOADER" = true ]; then
    echo "using loader"
    docker exec -i "${CONTAINER_NAME}" \
                loader \
                -f "${OUTPUT_FILE}" 
else
    echo "using load data infile"
    runSQL "LOAD DATA INFILE '${OUTPUT_FILE}'
            INTO TABLE ${TABLE_NAME}
            FIELDS TERMINATED BY ','
            LINES TERMINATED BY '\n'
            (ID, MSG, UPDATED_AT);"
fi

runSQL "SELECT COUNT(*) FROM ${TABLE_NAME}"

