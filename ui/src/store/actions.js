import Vue from "vue"
import tokenManager from "../helpers/token"

export default {
  updateHash: (context, hash, fields = 'md5|sha1|sha256|sha512|ssdeep|crc32|magic|size|exif|trid|tags|packer|first_submission|last_submission|last_scanned|submissions') => {
    context.commit('setHashContext', hash)
    context.dispatch('updateFileData', {
      hash: hash, fields: (fields instanceof Array) ? fields : fields.split('|')
    })
  },
  updateLoggedIn: (context, payload) => {
    if (!payload) {
      return
    }
    context.commit("setLoggedIn", Boolean(payload) && !tokenManager.isTokenExpired(payload))
    context.dispatch("updateUsername", payload)
  },
  updateUsername: (context, payload) => {
    if (!payload) {
      return
    }
    const jwtData = JSON.parse(atob(payload))
    context.dispatch('updateUserData', jwtData.name)
  },
  logOut: (context) => {
    Vue.$cookies.remove("JWTPayload")
    context.commit("setLoggedIn", false)
    context.commit("setUserData", {})
  },
  updateFileData: (context, {
    hash, fields = ['md5']
  }) => {
    console.log('fields', fields)
    Vue.prototype.$http.get(Vue.prototype.$api_endpoints.FILES + hash + '?' + ((fields) => {
      if (!fields || fields.length === 0) return '';
      return 'fields=' + fields.join(',');
    })(fields))
      .then((res) => {
        context.commit('setFileData', res.data)
      })
      .catch(() => {
        Vue.prototype.$awn.alert(
          "Sorry, we couldn't find the file you were looking for, please upload it to view the results!",
        )
        context.commit('setHashContext', "")
      })
  },
  updateUserData: (context, username) => {
    Vue.prototype.$http.get(Vue.prototype.$api_endpoints.USERS + username)
      .then((res) => {
        Vue.prototype.$http.get(Vue.prototype.$api_endpoints.USERS + username + "/avatar", {
          responseType: "arraybuffer",
        })
          .then((secRes) => {
            res.data.avatar = Buffer.from(secRes.data, 'binary').toString('base64')
            context.commit('setUserData', res.data)
          })
      })
      .catch(console.log)
  },
  addRemoveLike: (context, add) => {
    var data = context.getters.getUserData.likes
    if (add) {
      data.push(context.getters.getHashContext)
      context.commit('setLikes', data)
    } else {
      Vue._.pull(data, context.getters.getHashContext)
      context.commit('setLikes', data)
    }
  },
  updateComments: (context) => {
    Vue.prototype.$http.get(Vue.prototype.$api_endpoints.FILES + context.getters.getHashContext + "/?fields=comments")
      .then((res) => {
        context.commit('setComments', res.data.comments)
      })
      .catch(() => {
        Vue.prototype.$awn.alert("An Error Occured while updating the comments")
      })
  },
  updateLikes: (context) => {
    Vue.prototype.$http.get(Vue.prototype.$api_endpoints.USERS + context.getters.getUsername + "/?fields=likes")
      .then((res) => {
        context.commit('setLikes', res.data.likes)
      })
      .catch(() => {
        Vue.prototype.$awn.alert("An Error Occured while updating the likes")
      })
  },
  updateFollowing: (context) => {
    Vue.prototype.$http.get(Vue.prototype.$api_endpoints.USERS + context.getters.getUsername + "/?fields=following")
      .then((res) => {
        context.commit('setFollowing', res.data.following)
      })
      .catch((e) => {
        console.log(e)
        Vue.prototype.$awn.alert("An Error Occured while updating the followings")
      })
  },
}
