#!/bin/sh
# init-dev-server.sh

sleep 10

# Migration

# Starting server
# You can choose whether using air or using go to run the server
# go run ./cmd/server/main.go
go mod tidy
go mod vendor

air