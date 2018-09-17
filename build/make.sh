#!/bin/bash
# Build file for project. Will perform golint and gofmt before building.

sourceFiles="main.go"

gofmt -s -w .
golint ./...
go build ${sourceFiles}