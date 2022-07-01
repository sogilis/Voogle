#!/bin/bash

coverage-out() {
    go test -v -coverprofile coverage.out.tmp ./...
    cat coverage.out.tmp | grep -v "_dummy.go" > coverage.out
}

coverage(){
    coverage-out
    go tool cover -func=coverage.out
}

coverage-total() {
    coverage-out
    go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}'
}

coverage-html() {
    coverage-out
    go tool cover -html=coverage.out -o coverage.html
}

if [ -n "$1" ] && [ "total" == "$1" ]; then
    coverage-total
elif [ -n "$1" ] && [ "html" == "$1" ]; then
    coverage-html
else
    coverage
fi