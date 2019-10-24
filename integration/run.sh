#!/bin/bash

echo "# USER NOT FOUND"
curl http://localhost:8080/hello/john

echo "# ADD USER"
curl -XPUT http://localhost:8080/hello/john -H 'Content-type: application/json' -d '{"dateOfBirth": "2016-02-02"}'

echo "# USER FOUND"
curl http://localhost:8080/hello/john

echo "# UPDATE USER"
BIRTHDAY=$(date -v -1y -v +2d +%Y-%m-%d)
curl -XPUT http://localhost:8080/hello/john -H 'Content-type: application/json' -d '{"dateOfBirth": "'$BIRTHDAY'"}'

echo "# BIRTHDAY IS IN 2 DAYS"
curl http://localhost:8080/hello/john

echo "# UPDATE USER"
BIRTHDAY=$(date -v -1y -v -1d +%Y-%m-%d)
curl -XPUT http://localhost:8080/hello/john -H 'Content-type: application/json' -d '{"dateOfBirth": "'$BIRTHDAY'"}'

echo "# BIRTHDAY WAS YESTERDAY (in 364 days)"
curl http://localhost:8080/hello/john

echo "# SETTING TODAY AS BIRTHDAY"
BIRTHDAY=$(date -v -1y +%Y-%m-%d)
curl -XPUT http://localhost:8080/hello/john -H 'Content-type: application/json' -d '{"dateOfBirth": "'$BIRTHDAY'"}'

echo "# BIRTHDAY IS TODAY"
curl http://localhost:8080/hello/john
