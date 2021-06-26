#!/bin/sh
# run-prod.sh

docker-compose -f deployments/compose/docker-compose.yml up --build
