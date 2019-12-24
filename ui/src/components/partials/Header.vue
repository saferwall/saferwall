<template>
  <header class="dashboard-header">
    <router-link :to="this.$routes.HOME.path" class="logo">
      <img
        src="../../assets/imgs/logo-horizontal_rescaled.png"
        alt="Saferwall"
      />
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
    <nav class="dashboard-nav" :class="{ mobile: showinmobile }">
      <ul>
        <!-- <li><router-link to="/">Search</router-link></li> -->
        <li>
          <router-link :to="this.$routes.UPLOAD.path">
            <i class="icon fas fa-upload fa-2x"></i>
          </router-link>
        </li>
        <!-- <li><router-link to="/">Statistics</router-link></li> -->
        <li>
          <div class="profile">
            <div class="dropdown is-hoverable is-right">
              <div class="dropdown-trigger">
                <button
                  :class="getLoggedIn ? 'dropdown_text' : 'has-text-link'"
                  aria-haspopup="true"
                  aria-controls="dropdown-menu4"
                  @click="loginOrLogout"
                >
                  <span>{{ getLoggedIn ? getUsername : "Sign In" }}</span>
                </button>
              </div>
              <div
                class="dropdown-menu"
                id="dropdown-menu4"
                role="menu"
                v-if="getLoggedIn"
              >
                <div class="dropdown-content">
                  <div class="dropdown-item">
                    <div @click="loginOrLogout" class="has-text-danger">
                      <span>
                        <i class="icon fas fa-sign-out-alt fa-lg"></i>
                        Logout
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </li>
      </ul>
    </nav>
  </header>
</template>

<script>
import { mapGetters, mapActions } from "vuex"
export default {
  name: "Header",
  data() {
    return {
      hash: "",
      showinmobile: false,
    }
  },
  computed: {
    ...mapGetters(["getLoggedIn", "getUsername"]),
  },
  methods: {
    ...mapActions([
      "updateUsername",
      "updateLoggedIn",
      "logOut",
      "updateHash",
      "updateFileData",
    ]),
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
    searchByHash() {
      this.$http
        .get(`${this.$api_endpoints.FILES}${this.hash}/`, {
          validateStatus: (status) => status === 200,
        })
        .then((data) => {
          this.updateHash(this.hash)
          this.updateFileData(data)

          this.$router.push(this.$routes.SUMMARY.path + this.hash)
        })
        .catch(() => {
          this.$awn.alert(
            "Sorry, we couldn't find the file you were looking for, please upload it to view the results!",
          )
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
          padding-right: 10px;

          .has-text-danger {
            font-weight: 400;
            font-size: 14px;
            display: inline-block;

            &:hover {
              font-weight: 600;
            }
          }

          .dropdown-menu {
            min-width: 10rem;
          }
        }

        button {
          display: inline-block;
          width: 100%;
          color: rgba(black, 0.6);
          font-weight: 600;
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
