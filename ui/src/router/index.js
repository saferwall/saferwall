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
import {store} from "../store.js";

Vue.use(Router);

const storeState = store.state

let router = new Router({
  routes: [
    {
      path: "/",
      name: "Home",
      component: Home,
      meta: { title: "Home" }
    },
    {
      path: "/upload",
      name: "Upload",
      component: Upload,
      meta: {
        title: "Upload",
        requiresAuth: true
      }
    },
    {
      path: "/scanning",
      name: "Scanning",
      component: Scanning,
      meta: { title: "Scanning" }
    },
    {
      path: "/antivirus/:hash",
      name: "Antivirus",
      component: Antivirus,
      meta: { title: "Antivirus" }
    },
    {
      path: "/summary/:hash",
      name: "Summary",
      component: Summary,
      meta: { title: "Summary" }
    },
    {
      path: "/strings/:hash",
      name: "Strings",
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
      name: "Sign up",
      component: Signup,
      meta: { title: "Sign up", guest: true, layout: "unauthenticated" }
    }
  ]
});

router.beforeEach(function(to, from, next) {
  if (
    to.matched.some(record => record.meta.requiresAuth) &&
    from.path !== "/login"
  ) {
    if (!storeState.loggedIn) {
      next({
        name: "login",
        params: { nextUrl: to.fullPath }
      });
    } else {
      next();
    }
  } else if (to.matched.some(record => record.meta.guest)) {
    next();
  } else next();
});

export default router;
