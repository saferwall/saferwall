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
  },
  getLikes: state => {
    return state.userData.likes
  },
  getComments: state => {
    return state.fileData.data.comments
  },
  getAvatar: state =>{
    return state.userData.avatarBase64
  }
}