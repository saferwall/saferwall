// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from "vue"
import App from "./App"
import router from "./router"

// import layouts globally
import Default from "./layouts/Default.vue"
import Unauthenticated from "./layouts/Unauthenticated.vue"

import prodenv from '../config/prod.env'
import devenv from '../config/dev.env'

import store from './store/index'

import Vuelidate from "vuelidate"
import axios from "axios"
import VueCookies from "vue-cookies"


Vue.use(Vuelidate)
Vue.use(VueCookies)

// set default config
Vue.$cookies.config('7d')


Vue.prototype.$clipboard = (function(window, document, navigator) {
  var textArea, copy

  function isOS() {
    return navigator.userAgent.match(/ipad|iphone/i)
  }

  function createTextArea(text) {
    textArea = document.createElement("textArea")
    textArea.value = text
    document.body.appendChild(textArea)
  }

  function selectText() {
    var range, selection

    if (isOS()) {
      range = document.createRange()
      range.selectNodeContents(textArea)
      selection = window.getSelection()
      selection.removeAllRanges()
      selection.addRange(range)
      textArea.setSelectionRange(0, 999999)
    } else {
      textArea.select()
    }
  }

  function copyToClipboard() {
    document.execCommand("copy")
    document.body.removeChild(textArea)
  }

  copy = function(text) {
    createTextArea(text)
    selectText()
    copyToClipboard()
  }

  return {
    copy: copy,
  }
})(window, document, navigator)

Vue.component("default-layout", Default)
Vue.component("unauthenticated-layout", Unauthenticated)

Vue.directive("focus", {
  // When the bound element is inserted into the DOM...
  inserted: function(el) {
    // Focus the element
    el.focus()
  },
})

Vue.config.productionTip = false

let URL, API_ENDPOINTS;
console.log(process.env.NODE_ENV)
if(process.env.NODE_ENV === "development"){
  URL =  devenv.ROOT_API
  API_ENDPOINTS = devenv.API_ENDPOINTS
} else if (process.env.NODE_ENV === "production"){
  URL = prodenv.ROOT_API
  API_ENDPOINTS = devenv.API_ENDPOINTS
}

Vue.prototype.$http = axios.create({
  baseURL: URL,
  withCredentials: true,
})

Vue.prototype.$api_endpoints = API_ENDPOINTS


/* eslint-disable no-new */
new Vue({
  el: "#app",
  router,
  store,
  components: { App },
  render: (h) => h(App),
})
