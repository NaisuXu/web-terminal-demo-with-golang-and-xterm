#!/bin/sh
GOOS=linux GOARCH=amd64 go build -o webterminal_linux_x86-64
GOOS=linux GOARCH=arm GOARM=7 go build -o webterminal_linux_armv7
GOOS=linux GOARCH=arm GOARM=5 go build -o webterminal_linux_armv5
