'use strict'
module.exports = {
  NODE_ENV: '"production"',
  ROOT_API: "http://dev.api.saferwall.com",
  API_ENDPOINTS: {
    POST_FILE: "/v1/files/",
    GET_FILES: "/v1/files/",
    GET_USERS: "/v1/users/",
    POST_USER: "/v1/users/",
    AUTH_LOGIN: "/v1/auth/login/",
    AUTH_REGISTER: "/v1/users/",
    AUTH_CHANGE_PWD: "/v1/users/password/",
  },
  ROUTES: {
    HOME: {
      path: "/",
      name: "home",
      meta: {
        title: "home",
      }
    },
    UPLOAD: {
      path: "/upload/",
      name: "upload",
      meta: {
        title: "Upload",
        requiresAuth: true,
      },
    },
    SCANNING: {
      path: "/scanning/",
      name: "scanning",
      meta: {
        title: "Scanning"
      },
    },
    ANTIVIRUS: {
      path: "/antivirus/",
      name: "antivirus",
      meta: {
        title: "Antivirus"
      },
    },
    SUMMARY: {
      path: "/summary/",
      name: "summary",
      meta: {
        title: "Summary"
      },
    },
    STRINGS: {
      path: "/strings/",
      name: "strings",
      meta: {
        title: "Strings"
      },
    },
    LOGIN: {
      path: "/login/",
      name: "login",
      meta: {
        title: "Log in",
        guest: true,
        layout: "unauthenticated"
      },
    },
    SIGNUP: {
      path: "/signup/",
      name: "signUp",
      meta: {
        title: "Sign up",
        guest: true,
        layout: "unauthenticated"
      },
    },
    FORGOT_PWD: {
      path: "/forgot_password/",
      name: "forgotPassword",
      meta: {
        title: "Forgot Password?",
        guest: true,
        layout: "unauthenticated",
      },
    },
    RESET_PWD: {
      path: "/reset_password/",
      name: "resetPassword",
      meta: {
        title: "Reset Password",
        guest: true,
        layout: "unauthenticated",
      },
    }
  }
}
