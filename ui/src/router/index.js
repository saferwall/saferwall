import Vue from "vue";
import Router from "vue-router";
import Home from "@/components/pages/Home";
import Upload from "@/components/pages/Upload";
import Scanning from "@/components/pages/Scanning";
import Antivirus from "@/components/pages/Antivirus";
import Summary from "@/components/pages/Summary";
import Strings from "@/components/pages/Strings";
import Login from "@/components/pages/Login";
import Signup from "@/components/pages/Signup";
import ForgotPassword from "@/components/pages/ForgotPassword";
import ResetPassword from "@/components/pages/ResetPassword";
import { store } from "../store.js";

Vue.use(Router);

const storeLoggedIn = store.state.loggedIn;
const loadTokenFromCookie = () => {
  const token = Vue.cookie.get("JWTCookie");
  if (token !== null) {
    store.setLoggedIn(token);
    store.setUsername(token);
    return true;
  } else return false;
};

const isLogged = () => storeLoggedIn || loadTokenFromCookie();

const router = new Router({
  mode: "history",
  routes: [
    {
      path: "/",
      name: "home",
      component: Home,
      meta: { title: "Home" }
    },
    {
      path: "/upload",
      name: "upload",
      component: Upload,
      meta: {
        title: "Upload",
        requiresAuth: true
      }
    },
    {
      path: "/scanning",
      name: "scanning",
      component: Scanning,
      meta: { title: "Scanning" }
    },
    {
      path: "/antivirus/:hash",
      name: "antivirus",
      component: Antivirus,
      meta: { title: "Antivirus" }
    },
    {
      path: "/summary/:hash",
      name: "summary",
      component: Summary,
      meta: { title: "Summary" }
    },
    {
      path: "/strings/:hash",
      name: "strings",
      component: Strings,
      meta: { title: "Strings" }
    },
    {
      path: "/login",
      name: "login",
      component: Login,
      meta: { title: "Log in", guest: true, layout: "unauthenticated" }
    },
    {
      path: "/signup",
      name: "signUp",
      component: Signup,
      meta: { title: "Sign up", guest: true, layout: "unauthenticated" }
    },
    {
      path: "/forgot_password",
      name: "forgotPassword",
      component: ForgotPassword,
      meta: {
        title: "Forgot Password?",
        guest: true,
        layout: "unauthenticated"
      }
    },
    {
      path: "/reset_password",
      name: "resetPassword",
      component: ResetPassword,
      meta: {
        title: "Reset Password",
        guest: true,
        layout: "unauthenticated"
      }
    }
  ]
});

router.beforeEach(function(to, from, next) {
  if (to.matched.some(record => record.meta.requiresAuth) && !isLogged()) {
    next({
      name: "login",
      params: { nextUrl: to.fullPath }
    });
  } else if (to.matched.some(record => record.meta.guest)) {
    next();
  } else next();
});

export default router;
