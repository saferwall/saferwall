import Vue from 'vue'

export default {
  setHashContext(state, hash) {
    state.hashContext = hash
  },
  setLoggedIn(state, payload) {
    state.loggedIn = payload
  },
  setFileData(state, data) {
    state.fileData = data
  },
  setComments(state, data) {
    Vue.set(state.fileData.data, 'comments', data)
  },
  setUserData(state, data) {
    state.userData = data
  },
  setLikes(state, data) {
    Vue.set(state.userData, 'likes', data)
  },
  setFollowing(state, data) {
    state.userData.following = data
  }
}