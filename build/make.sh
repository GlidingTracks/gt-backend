#!/bin/bash
# Build file for project. Will perform golint and gofmt before building.

sourceFiles="main.go userHandler.go"

go fmt
golint
go build ${sourceFiles}