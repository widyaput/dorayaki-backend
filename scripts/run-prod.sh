#!/bin/sh
# run-prod.sh

docker-compose -d -f deployments/compose-prod/docker-compose.yml up --build
