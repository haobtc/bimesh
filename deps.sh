#!/bin/bash

export GOPATH=$PWD

function goget() {
    echo go get $1
    go get $1
    ecode=$?
    if [ $ecode -eq 0 ]; then
        echo done
    else
        exit $ecode
    fi
}

goget github.com/gorilla/websocket
goget github.com/bitly/go-simplejson
goget github.com/satori/go.uuid
goget gopkg.in/yaml.v2
goget github.com/stretchr/testify/assert
goget go.etcd.io/etcd/clientv3
