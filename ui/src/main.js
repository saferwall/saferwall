// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from "vue"
import App from "./App"
import router from "./router"

// import layouts globally
import Default from "./layouts/Default.vue"
import Unauthenticated from "./layouts/Unauthenticated.vue"

import endpoints from '../config/endpoints'
import Configuration from './helpers/config'

import store from './store/index'

import Vuelidate from "vuelidate"
import axios from "axios"
import VueCookies from "vue-cookies"
import VueAWN from "vue-awesome-notifications"
import VueLodash from 'vue-lodash'

// Google Analytics
import VueGtag from "vue-gtag";


require('../node_modules/vue-awesome-notifications/dist/styles/style.css')


const options = {
  position: "bottom-right",
  maxNotification: 5,
  animationDuration: 300,
  duration: {
    global: 20000,
  },
  minDurations: {
    "async-block": 1000,
    async: 1000,
  },
  labels : {
    info: "Information",
    success: "Success",
    warning: "Attention",
    alert: "Failure",
    async: "Processing",
  }
}


Vue.use(VueAWN, options)
Vue.use(VueLodash)
Vue.use(VueGtag, {
  config: { id: "UA-111524273-1" },
  appName: 'SaferWall',
}, router);
Vue.use(Vuelidate)
Vue.use(VueCookies)


Vue.prototype.$clipboard = (function (window, document, navigator) {
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

  copy = function (text) {
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
  inserted: function (el) {
    // Focus the element
    el.focus()
  },
})

Vue.config.productionTip = false


Vue.prototype.$http = axios.create({
  baseURL: Configuration.value('backendHost'),
  withCredentials: true,
})

Vue.prototype.$api_endpoints = endpoints

/* eslint-disable no-new */
new Vue({
  el: "#app",
  router,
  store,
  components: {
    App
  },
  render: (h) => h(App),
})
