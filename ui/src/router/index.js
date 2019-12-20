import Vue from "vue"
import Router from "vue-router"
import Home from "@/components/pages/Home"
import Upload from "@/components/pages/Upload"
import Scanning from "@/components/pages/Scanning"
import Antivirus from "@/components/pages/Antivirus"
import Summary from "@/components/pages/Summary"
import Strings from "@/components/pages/Strings"
import Login from "@/components/pages/Login"
import Signup from "@/components/pages/Signup"
import ForgotPassword from "@/components/pages/ForgotPassword"
import ResetPassword from "@/components/pages/ResetPassword"
import store from "../store/index"
import routes from '../../config/routes'


Vue.use(Router)

var ROUTES = routes
Vue.prototype.$routes = ROUTES

const loadTokenFromCookie = () => {
  const payload = Vue.$cookies.get("JWTPayload")
  if (payload !== null) {
    store.dispatch('setLoggedIn', payload)
    store.dispatch('setUsername', payload)
    return true
  } else return false
}

const isLogged = () => store.getters.getLoggedIn || loadTokenFromCookie()

const router = new Router({
  mode: "history",
  routes: [{
      path: ROUTES.HOME.path,
      name: ROUTES.HOME.name,
      component: Home,
      meta: {
        title: ROUTES.HOME.meta.title
      },
    },
    {
      path: ROUTES.UPLOAD.path,
      name: ROUTES.UPLOAD.name,
      component: Upload,
      meta: {
        title: ROUTES.UPLOAD.meta.title,
        requiresAuth: ROUTES.UPLOAD.meta.requiresAuth,
      },
    },
    {
      path: ROUTES.SCANNING.path,
      name: ROUTES.SCANNING.name,
      component: Scanning,
      meta: {
        title: ROUTES.SCANNING.meta.title
      },
    },
    {
      path: ROUTES.ANTIVIRUS.path + ":hash",
      name: ROUTES.ANTIVIRUS.name,
      component: Antivirus,
      meta: {
        title: ROUTES.ANTIVIRUS.meta.title
      },
    },
    {
      path: ROUTES.SUMMARY.path + ":hash",
      name: ROUTES.SUMMARY.name,
      component: Summary,
      meta: {
        title: ROUTES.SUMMARY.meta.title
      },
    },
    {
      path: ROUTES.STRINGS.path + ":hash",
      name: ROUTES.STRINGS.name,
      component: Strings,
      meta: {
        title: ROUTES.STRINGS.meta.title
      },
    },
    {
      path: ROUTES.LOGIN.path,
      name: ROUTES.LOGIN.name,
      component: Login,
      meta: {
        title: ROUTES.LOGIN.meta.title,
        guest: ROUTES.LOGIN.meta.guest,
        layout: ROUTES.LOGIN.meta.layout
      },
    },
    {
      path: ROUTES.SIGNUP.path,
      name: ROUTES.SIGNUP.name,
      component: Signup,
      meta: {
        title: ROUTES.SIGNUP.meta.title,
        guest: ROUTES.SIGNUP.meta.guest,
        layout: ROUTES.SIGNUP.meta.layout
      },
    },
    {
      path: ROUTES.FORGOT_PWD.path,
      name: ROUTES.FORGOT_PWD.name,
      component: ForgotPassword,
      meta: {
        title: ROUTES.FORGOT_PWD.meta.title,
        guest: ROUTES.FORGOT_PWD.meta.guest,
        layout: ROUTES.FORGOT_PWD.meta.layout,
      },
    },
    {
      path: ROUTES.RESET_PWD.path,
      name: ROUTES.RESET_PWD.name,
      component: ResetPassword,
      meta: {
        title: ROUTES.RESET_PWD.meta.title,
        guest: ROUTES.RESET_PWD.meta.guest,
        layout: ROUTES.RESET_PWD.meta.layout,
      },
    },
  ],
})

router.beforeEach(function (to, from, next) {
  if (to.matched.some((record) => record.meta.requiresAuth) && !isLogged()) {
    next({
      name: ROUTES.LOGIN.name,
      params: {
        nextUrl: to.fullPath
      },
    })
  } else if (to.matched.some((record) => record.meta.guest)) {
    next()
  } else next()
})

export default router
