import Vue from "vue"
import isTokenExpired from "../helpers/token"

export default {
  updateHash: (context, hash) => {
    context.commit('setHashContext', hash)
  },
  updateLoggedIn: (context, payload) => {
    if (!payload) {
      return
    }
    context.commit("setLoggedIn", Boolean(payload) && !isTokenExpired(payload))
  },
  updateUsername: (context, payload) => {
    if (!payload) {
      return
    }
    const jwtData = JSON.parse(atob(payload))
    context.commit("setUsername", jwtData.name || null)
  },
  logOut: (context) => {
    Vue.$cookies.remove("JWTPayload")
    context.commit("setLoggedIn", false)
    context.commit("setUsername", "")
  },
  updateFileData : (context, payload) => {
    context.commit('setFileData', payload)
  }
}
