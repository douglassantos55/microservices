#!/usr/bin/env sh

protoc --go_out=. --go_opt=paths=source_relative proto/$1.proto
