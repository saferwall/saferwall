export default {
  setHashContext(state, hash){
    state.hashContext = hash
  },
  setLoggedIn(state, payload){
    state.loggedIn = payload
  },
  setUsername(state, username){
    state.username = username
  },
  getTokenExpirationDate(state,payload){
    const jwtData = JSON.parse(atob(payload))
    if (!jwtData.exp) {
      return null
    }
  
    const date = new Date(0)
    date.setUTCSeconds(jwtData.exp)
  
    return date
  },
  isTokenExpired(state, payload){
    const expirationDate = state.actions.getTokenExpirationDate(payload)
    return expirationDate < new Date()
  }
}
