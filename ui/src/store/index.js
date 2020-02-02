import Vuex from "vuex"
import Vue from "vue"

import globalGetters from "./getters"
import globalMutations from "./mutations"
import globalActions from "./actions"

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    hashContext: "",
    fileData : {},
    loggedIn: false,
    userData: {}
  },
  getters: globalGetters,
  mutations: globalMutations,
  actions: globalActions

})
