#!/usr/bin/env sh

rm -f $1/proto/$2.pb.go
protoc --go_out=. --go_opt=paths=source_relative $1/proto/$2.proto
