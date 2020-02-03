export default {
  getHashContext: state => {
    return state.hashContext
  },
  getLoggedIn: state => {
    return state.loggedIn
  },
  getUsername: state => {
    return state.userData.username
  },
  getFileData: state => {
    return state.fileData
  },
  getUserData: state => {
    return state.userData
  }
}