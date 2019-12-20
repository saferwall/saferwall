<template>
  <header class="dashboard-header">
    <router-link :to="this.$routes.HOME.path" class="logo">
      <img src="../../assets/imgs/logo-horizontal.png" alt="" />
    </router-link>
    <div class="mobile-nav" @click="showinmobile = !showinmobile">
      <i class="icon ion-android-menu"></i>
    </div>
    <div class="header-search" :class="{ active: showinmobile }">
      <input
        type="search"
        placeholder="Quick lookup file hash, URL or IP."
        v-model="hash"
        @keyup.enter="searchByHash"
      />
      <button type="submit" @click.prevent="searchByHash">
        <i class="icon ion-ios-search"></i>
      </button>
    </div>
    <transition name="slide-fade">
      <notification
        :style="notificationStyling"
        type="is-danger"
        @closeNotif="close()"
        v-if="notifActive"
      >
        {{ notificationError }}
      </notification>
    </transition>
    <nav class="dashboard-nav" :class="{ mobile: showinmobile }">
      <ul>
        <!-- <li><router-link to="/">Search</router-link></li> -->
        <li>
          <router-link :to="this.$routes.UPLOAD.path">
            <i class="icon fas fa-upload fa-2x"></i>
          </router-link>
        </li>
        <!-- <li><router-link to="/">Statistics</router-link></li> -->
        <li class="has-dropdown" @click="dropdownActive = !dropdownActive">
          <div class="profile">
            <span>{{ getUsername || "" }}</span>
            <img src="../../assets/imgs/avatar.jpg" alt="" />
          </div>
          <ul class="dropdown-container" :class="{ active: dropdownActive }">
            <li>
              <router-link to="/">
                <i class="icon ion-grid"></i>
                Your Submissions
              </router-link>
            </li>
            <li>
              <router-link to="/">
                <i class="icon ion-android-settings"></i>
                Settings
              </router-link>
            </li>
            <li>
              <button
                :class="getLoggedIn ? 'has-text-danger' : 'has-text-link'"
                @click="loginOrLogout"
              >
                <i
                  class="icon"
                  :class="getLoggedIn ? 'ion-log-out' : 'ion-log-in'"
                ></i>
                {{ getLoggedIn ? "Sign Out" : "Sign In" }}
              </button>
            </li>
          </ul>
        </li>
      </ul>
    </nav>
  </header>
</template>

<script>
import { mapGetters, mapActions } from "vuex"
import Notification from "@/components/elements/Notification"
export default {
  name: "Header",
  data() {
    return {
      hash: "",
      notificationError: "",
      notifActive: false,
      dropdownActive: false,
      showinmobile: false,
      notificationStyling: {
        width: "fit-content",
        height: "auto",
        position: "absolute",
        left: "50%",
        top: "10px",
        transform: "translateX(-50%)",
      },
    }
  },
  computed: {
    ...mapGetters(["getLoggedIn", "getUsername"]),
  },
  components: {
    notification: Notification,
  },
  methods: {
    ...mapActions(["updateUsername", "updateLoggedIn", "logOut", "updateHash"]),
    showMobileSearch() {},
    loginOrLogout() {
      if (this.getLoggedIn) {
        this.logOut()
        this.$router.go(this.$routes.HOME.path)
      } else {
        this.$router.push(this.$routes.LOGIN.path)
      }
    },

    getJWTPayload() {
      const payload = this.$cookies.get("JWTPayload")
      return payload
    },
    close() {
      this.notifActive = false
    },
    searchByHash() {
      this.$http
        .get(`${this.$api_endpoints.FILES}${this.hash}/`, {
          validateStatus: (status) => status === 200,
        })
        .then(() => {
          this.updateHash(this.hash)
          this.$router.push(this.$routes.SUMMARY.path + this.hash)
        })
        .catch(() => {
          this.notifActive = true
          this.notificationError =
            "Sorry, we couldn't find the file you were looking for, please upload it to view the results!"
          setTimeout(() => {
            this.notifActive = false
          }, 3000)
        })
    },
  },

  mounted() {
    const payload = this.getJWTPayload()
    this.updateLoggedIn(payload)
    this.updateUsername(payload)
  },
}
</script>

<style scoped lang="scss">
@import "../../assets/scss/variables";
$header-height: 50px;

header.dashboard-header {
  margin-right: 20px;
  background: #fff;
  height: $header-height;
  line-height: $header-height;
  box-shadow: 0 2px 3px rgba(10, 10, 10, 0.1), 0 0 0 1px rgba(10, 10, 10, 0.1);
  z-index: 99;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  display: flex;

  .logo {
    float: left;
    line-height: 62px;
    height: $header-height;
    width: 200px;
    text-align: center;

    img {
      height: calc(#{$header-height - 10px});
      display: inline-block;
      margin-top: 5px;
    }
  }

  .mobile-nav {
    // margin-left:20px;
    float: right;
    height: $header-height;
    width: 50px;
    line-height: $header-height;
    border-left: solid 1px rgba(10, 10, 10, 0.1);
    border-right: solid 1px rgba(10, 10, 10, 0.1);
    text-align: center;

    @media screen and (min-width: 792px) {
      display: none;
    }
  }

  .header-search {
    display: inline-block;
    width: 500px;
    height: $header-height;
    position: relative;
    border-left: solid 1px rgba(black, 0.1);
    border-right: solid 1px rgba(black, 0.1);
    flex: 1;

    @media screen and (max-width: 1086px) {
      width: 400px;
    }

    @media screen and (max-width: 997px) {
      width: 300px;
    }

    @media screen and (max-width: 890px) {
      width: 200px;
    }

    input {
      float: left;
      height: $header-height;
      width: calc(100% - #{$header-height});
      padding: 0 10px;
      border: 0;
      font-size: 13px;
    }

    button {
      float: right;
      width: 50px;
      height: $header-height;
      line-height: $header-height;
      text-align: center;
      color: $primary-color;
      border: 0;
      background-color: transparent;
      font-size: 24px;
      cursor: pointer;
    }
  }

  nav.dashboard-nav {
    float: right;
    &,
    * {
      transition: all 0s;
    }

    @media screen and (max-width: 792px) {
      display: none;
    }

    .icon {
      display: inline-block;
    }
    ul {
      li {
        display: inline-block;
        line-height: $header-height;

        a {
          display: inline-block;
          margin: 0 15px 0 15px;
          font-size: 14px;
          color: #2c3e50;
          font-weight: 500;
          transition: all 0.2s;

          &:hover {
            color: $primary-color;
          }
        }

        & > .profile {
          cursor: pointer;
          border-left: solid 1px rgba(10, 10, 10, 0.1);
          padding-left: 10px;

          span {
            font-size: 14px;
            font-weight: 500;
            margin-right: 10px;
          }

          img {
            display: inline-block;
            vertical-align: middle;
            height: 40px;
            width: 40px;
            border-radius: 50%;
            border: solid 1px rgba(10, 10, 10, 0.1);
          }
        }

        &.has-dropdown {
          position: relative;

          .dropdown-container {
            transform: scale(0);
            transform-origin: 100% 0;
            transition: all 0.2s;
            position: absolute;
            top: 40px;
            right: 0;
            width: 200px;
            padding: 0 10px;
            background: #fff;
            // box-shadow: 0 2px 3px rgba(10, 10, 10, 0.1), 0 0 0 1px rgba(10, 10, 10, 0.1);
            box-shadow: -2px 3px 10px rgba(10, 10, 10, 0.1);
            border-radius: 4px;

            &.active {
              transform: scale(1);
            }

            li {
              line-height: 30px;
              display: inline;

              a,
              button {
                display: inline-block;
                width: 100%;
                color: rgba(black, 0.6);
                font-weight: 500;
                font-size: 14px;

                &:hover {
                  color: $primary-color !important;
                }
              }

              button {
                text-align: left;
                padding: 0 5px 0 5px;
                transition: all 0.2s;
                background: none;
                border: none;
                cursor: pointer;
              }
            }
          }
        }
      }
    }

    &.mobile {
      display: block;
      position: fixed;
      top: 100px;
      width: 100%;
      background-color: #fff;
      z-index: 9999999999;
      box-shadow: 0 5px 10px rgba(10, 10, 10, 0.1);
      padding-bottom: 10px;

      li {
        width: 100%;
        text-align: center;

        & > .profile {
          border-left: 0;
        }
      }

      ul.dropdown-container {
        right: calc(50% - 100px) !important;

        li a {
          text-align: left;
        }
      }
    }
  }

  .slide-fade-enter-active {
    transition: all 0.3s ease;
  }
  .slide-fade-leave-active {
    transition: all 0.8s cubic-bezier(1, 0.5, 0.8, 1);
  }
  .slide-fade-enter,
  .slide-fade-leave-to {
    transform: translateX(10px);
    opacity: 0;
  }

  @media screen and (max-width: 792px) {
    padding-right: 20px;

    .header-search {
      position: fixed;
      top: 50px;
      left: 0;
      width: 100%;
      background-color: #fff;
      margin-left: 0;
      height: 0;
      padding: 0;
      overflow: hidden;

      &.active {
        padding: 0 10px;
        height: 50px;
      }
    }

    nav.dashboard-nav {
      ul {
        li {
          & > .profile {
            span {
              display: none;
            }
          }
        }
      }
    }
  }
}
</style>
