'use strict'

const mesh = require('./lib/mesh')

mesh.start()
mesh.on('ready', () => {
  onReady().then(console.log, console.error)
})

async function onReady() {
  let r = await mesh.call('calc::add', [7, 9])
  console.info(r)
}
