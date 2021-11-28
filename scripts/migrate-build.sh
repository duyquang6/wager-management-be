#!/bin/bash
REPO="github.com/duyquang6/wager-management-be"
NOW=$(date +'%Y-%m-%d_%T')

go build -ldflags "-X $REPO/internal/buildinfo.buildID=`git rev-parse --short HEAD` -X $REPO/internal/buildinfo.buildTime=$NOW" -o bin/migrate cmd/migrate/main.go
