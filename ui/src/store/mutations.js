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
  setFileData(state, data){
    state.fileData = data
  }
}
