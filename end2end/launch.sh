#!/usr/bin/env sh
USER_AUTH=`cat ../.env | grep USER_AUTH | awk -F'=' '{print $2}' | tr -d '\n'`
PWD_AUTH=`cat ../.env | grep PWD_AUTH | awk -F'=' '{print $2}' | tr -d '\n'`

echo E2E_API_ENDPOINT=http://localhost:4444 E2E_USER_NAME=$USER_AUTH E2E_USER_PWD=$PWD_AUTH go test ./... -v

E2E_API_ENDPOINT=http://localhost:4444 E2E_USER_NAME=$USER_AUTH E2E_USER_PWD=$PWD_AUTH go test ./... -v
