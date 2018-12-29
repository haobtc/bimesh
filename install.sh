#!/bin/bash

export GOPATH=$PWD

echo testing jsonrpc
go test jsonrpc

go install bimeshd


