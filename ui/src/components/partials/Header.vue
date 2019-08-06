<template>
  <header class="dashboard-header">
    <router-link to="/" class="logo">
      <img src="../../assets/imgs/logo.png" alt="" />
    </router-link>
    <div class="mobile-nav" @click="showinmobile = !showinmobile">
      <i class="icon ion-android-menu"></i>
    </div>
    <div class="header-search" :class="{ active: showinmobile }">
      <input type="text" placeholder="Quick lookup file hash, URL or IP." />
      <button type="submit">
        <i class="icon ion-ios-cloud-upload-outline"></i>
      </button>
    </div>
    <nav class="dashboard-nav" :class="{ mobile: showinmobile }">
      <ul>
        <li><router-link to="/">Search</router-link></li>
        <li><router-link to="/">Upload</router-link></li>
        <li><router-link to="/">Statistics</router-link></li>
        <li class="has-dropdown" @click="dropdownActive = !dropdownActive">
          <div class="profile">
            <span>{{ storeState.username || "" }}</span>
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
                :class="
                  storeState.loggedIn ? 'has-text-danger' : 'has-text-link'
                "
                @click="loginOrLogout"
              >
                <i
                  class="icon"
                  :class="storeState.loggedIn ? 'ion-log-out' : 'ion-log-in'"
                ></i>
                {{ storeState.loggedIn ? "Sign Out" : "Sign In" }}
              </button>
            </li>
          </ul>
        </li>
      </ul>
    </nav>
  </header>
</template>

<script>
import { store } from "../../store.js";
export default {
  name: "Header",
  data() {
    return {
      dropdownActive: false,
      showinmobile: false,
      storeState: store.state
    };
  },
  methods: {
    showMobileSearch() {},
    loginOrLogout() {
      if (this.storeState.loggedIn) {
        store.logOut();
        this.$router.go('/')
      } else {
        this.$router.push("/login");
      }
    },

    getJWTToken() {
      const token = this.$cookie.get("JWTCookie");
      return token;
    }
  },

  created() {
    const token = this.getJWTToken()
    store.setLoggedIn(token);
    store.setUsername(token);
  }
};
</script>

<style scoped lang="scss">
@import "../../assets/scss/variables";

header.dashboard-header {
  padding-right: 20px;
  background: #fff;
  height: 50px;
  line-height: 50px;
  box-shadow: 0 2px 3px rgba(10, 10, 10, 0.1), 0 0 0 1px rgba(10, 10, 10, 0.1);
  z-index: 99;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;

  .logo {
    float: left;
    line-height: 62px;
    height: 50px;
    width: 200px;
    text-align: center;

    img {
      height: 20px;
    }
  }

  .mobile-nav {
    // margin-left:20px;
    float: right;
    height: 50px;
    width: 50px;
    line-height: 50px;
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
    height: 50px;
    position: relative;
    border-left: solid 1px rgba(black, 0.1);
    border-right: solid 1px rgba(black, 0.1);

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
      height: 50px;
      width: calc(100% - 50px);
      padding: 0 10px;
      border: 0;
      font-size: 13px;
    }

    button {
      float: right;
      width: 50px;
      height: 50px;
      line-height: 50px;
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

    ul {
      li {
        display: inline-block;
        line-height: 50px;

        a {
          display: inline-block;
          padding: 0 5px;
          font-size: 14px;
          color: #2c3e50;
          font-weight: 500;
          transition: all 0.2s;

          &:hover {
            color: $primary-color;
          }
        }

        & > .profile {
          margin-left: 10px;
          cursor: pointer;
          border-left: solid 1px rgba(10, 10, 10, 0.1);
          padding-left: 10px;

          span {
            font-size: 14px;
            padding-right: 10px;
            font-weight: 500;
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
