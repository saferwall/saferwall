// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from "vue"
import App from "./App"
import router from "./router"

// import layouts globally
import Default from "./layouts/Default.vue"
import Unauthenticated from "./layouts/Unauthenticated.vue"

import prodenv from '../config/prod.env'

import Vuelidate from "vuelidate"

import axios from "axios"

import VueCookie from "vue-cookie"

Vue.use(Vuelidate)
Vue.use(VueCookie)

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

let URL;
if(process.env.NODE_ENV === "development"){
  URL =  prodenv.ROOT_API
} else{
  URL = prodenv.ROOT_API
}

Vue.prototype.$http = axios.create({
  baseURL: URL,
})
// alert(URL)
/* eslint-disable no-new */
new Vue({
  el: "#app",
  router,
  components: { App },
  render: (h) => h(App),
})
