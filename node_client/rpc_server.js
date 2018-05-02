'use strict'

const mesh = require('./lib/mesh')

mesh.start('ws://localhost:8080/jsonrpc/ws')
mesh.register('calc')

mesh.on('calc::add', function(msg) {
  let ret = msg.params[0] + msg.params[1]
  mesh.result(msg.id, ret)
})
