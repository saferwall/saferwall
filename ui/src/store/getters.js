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
    if (Object.entries(state.fileData).length === 0 && state.fileData.constructor === Object) return []
    return state.fileData.data.comments ? state.fileData.data.comments : []
  },
  getNbComments: state => {
    if (Object.entries(state.fileData).length === 0 && state.fileData.constructor === Object) return 0
    return state.fileData.data.comments ? state.fileData.data.comments.length : 0
  },
  getAvatar: state => {
    return state.userData.avatarBase64
  }
}
