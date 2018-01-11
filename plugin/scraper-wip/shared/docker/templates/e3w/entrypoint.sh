#!/bin/sh
glide install --strip-vendor
go run main.go -conf /data/conf.d/e3w/config.ini -front-dir /go/src/github.com/roscopecoltran/e3w/static/dist