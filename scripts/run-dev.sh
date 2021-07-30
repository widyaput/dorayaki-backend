#!/bin/sh
# run-prod.sh

docker-compose -f deployments/compose/docker-compose.yml -p dorayaki up --build
