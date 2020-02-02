export default {
  setHashContext(state, hash){
    state.hashContext = hash
  },
  setLoggedIn(state, payload){
    state.loggedIn = payload
  },
  setFileData(state, data){
    state.fileData = data
  },
  setUserData(state, data){
    state.userData = data
  }
}
