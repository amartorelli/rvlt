#!/bin/bash

docker-compose -f ./integration/docker-compose.yml up -d  --force-recreate

echo "# USER NOT FOUND"
curl http://localhost:8080/hello/john

echo "# ADD USER"
curl -XPUT http://localhost:8080/hello/john -H 'Content-type: application/json' -d '{"dateOfBirth": "2016-02-02"}'

echo "# USER FOUND"
curl http://localhost:8080/hello/john

echo "# UPDATE USER"
BIRTHDAY=$(date -v -1y -v +2d +%Y-%m-%d)
curl -XPUT http://localhost:8080/hello/john -H 'Content-type: application/json' -d '{"dateOfBirth": "'$BIRTHDAY'"}'

echo "# NEW BIRTHDAY IN 2 DAYS"
curl http://localhost:8080/hello/john

echo "# UPDATE USER"
BIRTHDAY=$(date -v -1y -v -1d +%Y-%m-%d)
curl -XPUT http://localhost:8080/hello/john -H 'Content-type: application/json' -d '{"dateOfBirth": "'$BIRTHDAY'"}'

echo "# NEW BIRTHDAY YESTERDAY DAYS (in 364 days)"
curl http://localhost:8080/hello/john

echo "# SETTING TODAY AS BIRTHDAY"
BIRTHDAY=$(date -v -1y +%Y-%m-%d)
curl -XPUT http://localhost:8080/hello/john -H 'Content-type: application/json' -d '{"dateOfBirth": "'$BIRTHDAY'"}'

echo "# BIRTHDAY IS TODAY"
curl http://localhost:8080/hello/john

# docker-compose -f ./integration/docker-compose.yml down
