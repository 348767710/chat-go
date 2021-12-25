#!/bin/bash
export GOOS=linux GOARCH=amd64
# export GOPATH=$GOPATH:/Users/pk/Desktop/project/
go build -ldflags "-w -s" -o chat-goim main.go && chmod a+x chat-goim
