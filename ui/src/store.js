import Vue from "vue"

const getTokenExpirationDate = (token) => {
  const jwtData = JSON.parse(atob(token.split(".")[1]))
  if (!jwtData.exp) {
    return null
  }

  const date = new Date(0)
  date.setUTCSeconds(jwtData.exp)

  return date
}

const isTokenExpired = (token) => {
  const expirationDate = getTokenExpirationDate(token)
  return expirationDate < new Date()
}

export const store = {
  debug: true,
  state: {
    loggedIn: false,
    username: "",
  },
  // call setters on mounted/created of the header
  setLoggedIn(token) {
    if (!token) {
      return
    }
    this.state.loggedIn = Boolean(token) && !isTokenExpired(token)
  },
  setUsername(token) {
    if (!token) {
      return
    }
    const jwtData = JSON.parse(atob(token.split(".")[1]))
    this.state.username = jwtData.name || null
  },
  logOut() {
    Vue.cookie.delete("JWTCookie")
    this.state.loggedIn = false
  },
}
