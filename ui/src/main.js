// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'

// import layouts globally
import Default from './layouts/Default.vue'
import Unauthenticated from './layouts/Unauthenticated.vue'
import Error from './layouts/Error.vue'

import Vuelidate from 'vuelidate'
Vue.use(Vuelidate)

import VueCookie from 'vue-cookie'
Vue.use(VueCookie)

Vue.component('default-layout', Default)
Vue.component('unauthenticated-layout', Unauthenticated)
Vue.component('error-layout', Error)

Vue.directive('focus', {
  // When the bound element is inserted into the DOM...
  inserted: function (el) {
    // Focus the element
    el.focus()
  }
})

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})
