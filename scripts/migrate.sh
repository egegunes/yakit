#!/bin/bash

DBCONTAINERNAME=yakitdb
DBNAME=yakit
DBUSER=yakit
MIGRATIONSFOLDER=$(pwd)/db/migrations

function migrate {
    name=$1

    migration="$MIGRATIONSFOLDER/$name"

    if [[ ! -f "$migration" ]]; then
        echo "$migration does not exist"
        exit 1
    fi

    echo "Running $name"

    sudo docker cp "$migration" "$DBCONTAINERNAME":/"$name"
    sudo docker exec "$DBCONTAINERNAME" psql -U "$DBUSER" -f "/$name" "$DBNAME"

    exit 0
}

migrate $1
