import Vue from "vue"

const getTokenExpirationDate = (payload) => {
  const jwtData = JSON.parse(atob(payload))
  if (!jwtData.exp) {
    return null
  }

  const date = new Date(0)
  date.setUTCSeconds(jwtData.exp)

  return date
}

const isTokenExpired = (payloda) => {
  const expirationDate = getTokenExpirationDate(payloda)
  return expirationDate < new Date()
}

export const store = {
  debug: true,
  state: {
    loggedIn: false,
    username: "",
    hash: ""
  },
  // call setters on mounted/created of the header
  setLoggedIn(payload) {
    if (!payload) {
      return
    }
    this.state.loggedIn = Boolean(payload) && !isTokenExpired(payload)
  },
  setUsername(payload) {
    if (!payload) {
      return
    }
    const jwtData = JSON.parse(atob(payload))
    this.state.username = jwtData.name || null
  },
  logOut() {
    Vue.$cookies.remove("JWTPayload")
    this.state.loggedIn = false
  },
  setHash(hash) {
    this.state.hash = hash
  }
}
