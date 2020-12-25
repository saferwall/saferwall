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
    return state.fileData.comments ? state.fileData.comments : []
  },
  getNbComments: state => {
    if (Object.entries(state.fileData).length === 0 && state.fileData.constructor === Object) return 0
    return state.fileData.comments ? state.fileData.comments.length : 0
  },
  getAvatar: state => {
    return state.userData.avatar
  },
  getFollowing: state => {
    return state.userData.following ? state.userData.following : []
  },
  isPE: state => {
    if (!state.fileData || !state.fileData.magic || Object.entries(state.fileData).length === 0)
      return false
    return state.fileData.magic.substring(0, 2) === "PE"
  }
}
