#!/usr/bin/env sh

cd "${0%/*}"

export GO111MODULE=on
go mod vendor

export GO111MODULE=off

gopherjs build main.go -o main.js
