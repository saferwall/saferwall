export default {
  getHashContext: state => {
    return state.hashContext
  },
  getLoggedIn: state => {
    return state.loggedIn
  },
  getUsername: state => {
    return state.username
  },
  getFileData: state => {
    return state.fileData
  }
}