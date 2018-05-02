#!/bin/bash

export GOPATH=$PWD

echo go get github.com/gorilla/websocket
go get github.com/gorilla/websocket

echo go get github.com/bitly/go-simplejson
go get github.com/bitly/go-simplejson

echo go get github.com/satori/go.uuid
go get github.com/satori/go.uuid

echo go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/assert

echo testing mesh
go test mesh

go install bimeshd


