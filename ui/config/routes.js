module.exports = {
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
    COMMENTS: {
      path: "/comments/",
      name: "comments",
      meta: {
        title: "Comments"
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
      path: "/forgot-password/",
      name: "forgotPassword",
      meta: {
        title: "Forgot Password?",
        guest: true,
        layout: "unauthenticated",
      },
    },
    RESET_PWD: {
      path: "/reset-password/",
      name: "resetPassword",
      meta: {
        title: "Reset Password",
        guest: true,
        layout: "unauthenticated",
      },
    },
    SETTINGS: {
      path: "/settings/",
      name: "settings",
      meta: {
        title: "Settings",
        requiresAuth: true,
      },
    },
  }