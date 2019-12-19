'use strict'
const merge = require('webpack-merge')
const prodEnv = require('./prod.env')

module.exports = {
  NODE_ENV: '"development"', // can be removed, webpack 4 automatically makes process.env.NODE_ENV available in source code
  ROOT_API: "http://dev.api.saferwall.com",
  API_ENDPOINTS: {
    FILES: "/v1/files/",
    USERS: "/v1/users/",
    AUTH_LOGIN: "/v1/auth/login/",
    AUTH_REGISTER: "/v1/users/",
    AUTH_CHANGE_PWD: "/v1/users/password/",
  },
}
