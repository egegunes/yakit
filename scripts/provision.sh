#!/bin/bash

DBCONTAINERNAME=yakitdb
DBNAME=yakit
DBUSER=yakit
MIGRATIONSFOLDER=$(pwd)/db/migrations
DUMMYDATAFOLDER=$(pwd)/db/dummy

function run_migrations {
    for migration in "$MIGRATIONSFOLDER"/*.sql; do
        name=$(basename "$migration")
        sudo docker cp "$migration" "$DBCONTAINERNAME":/"$name"
        sudo docker exec "$DBCONTAINERNAME" psql -U "$DBUSER" -f "/$name" "$DBNAME"
    done
}


function load_dummy_data {
    for data in $DUMMYDATAFOLDER/*.sql; do
        name=$(basename $data)
        sudo docker cp "$data" "$DBCONTAINERNAME":/"$name"
        sudo docker exec "$DBCONTAINERNAME" psql -U "$DBUSER" -f "/$name" "$DBNAME"
    done
}

function main {
    run_migrations
    load_dummy_data
}

main
