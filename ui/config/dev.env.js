'use strict'
const merge = require('webpack-merge')
const prodEnv = require('./prod.env')

module.exports = merge(prodEnv, {
  NODE_ENV: '"development"', // can be removed, webpack 4 automatically makes process.env.NODE_ENV available in source code
  ROOT_API: "http://dev.api.saferwall.com"

})
