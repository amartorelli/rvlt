#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "helloworld" --dbname "helloworld" <<-EOSQL
    CREATE TABLE birthdays(
        username VARCHAR(64) PRIMARY KEY,
        birthday DATE NOT NULL
    );
EOSQL
