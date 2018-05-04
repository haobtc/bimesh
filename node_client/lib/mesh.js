'use strict'

const EventEmitter = require('events')
const WebSocket = require('ws')

class MeshContext extends EventEmitter
{
  constructor() {
    super()
    this.counter = 1;
    this.rpcCallbacks = {}
    this.ws = null
    this.connId = null
    this.pendingServiceNames = []
  }

  start(meshServer) {
    meshServer = meshServer || process.env.BIMESH_CONNECT
    var self = this
    this.rpcCallbacks = {}
    this.connId = null
    this.ws = new WebSocket(meshServer)
    this.ws.on('open', () => {
      self.call('builtin-services::getId', [], (err, ret) => {
        if(err) {
          throw err
        } else {
          this.connId = ret
        }

        if(self.pendingServiceNames.length > 0) {
          self.register.apply(self, self.pendingServiceNames)
          self.pendingServiceNames = []
        }
        self.emit('ready', () => {})
      })
    })

    this.ws.on('close', () =>{
      self.ws = null
      self.connId = null
    })

    this.ws.on('message', (data) => {
      let msg = JSON.parse(data)
      if(msg.id && !msg.method) {
        //msg is return or error
        let cb = this.rpcCallbacks[msg.id]
        if(cb) {
          delete this.rpcCallbacks[msg.id]
          cb(msg.error, msg.result)
        }
      } else if(msg.method) {
        // msg is request/notify
        this.emit(msg.method, msg)
      } else {
        // TODO: assert error
      }
    })
  }

  sendJSON(data) {
    this.ws.send(JSON.stringify(data))
  }

  call(serviceMethod, params, cb) {
    if(!cb) {
      return this.callAsync(serviceMethod, params)
    }
    let callId = this.counter++
    this.rpcCallbacks[callId] = cb
    this.sendJSON({
      id: callId,
      method: serviceMethod,
      params: params
    })
    return callId
  }

  callAsync(serviceMethod, params) {
    params = params || []
    let self = this
    return new Promise(
      (resolve, reject) => {
        this.call(serviceMethod, params, (err, res) => {
          if(err) {
            reject(err)
          } else {
            resolve(res)
          }
        })
      })
  }

  result(msgId, result) {
    this.sendJSON({
      id: msgId,
      result: result
    })
  }

  register() {
    let names = []
    for(let name of arguments) {
      names.push(name)
    }
    if(names.length < 1) {
      throw new Error('must add service 1 or more names')
    }

    if(this.ws && this.connId) {
      this.call(
        'builtin-services::register', names, (err, ret) => {
          if(err) {
            throw err
          }
          // ret
        })
    } else {
      this.pendingServiceNames = this.pendingServiceNames.concat(names)
    }
  }
}

module.exports = new MeshContext()
