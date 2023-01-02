#!/usr/bin/env sh

rm -f proto/$1.pb.go
protoc --go_out=. --go_opt=paths=source_relative proto/$1.proto
