<template>
  <header class="dashboard-header">
    <router-link :to="this.$routes.HOME.path" class="logo">
      <img
        src="../../assets/imgs/logo-horizontal_rescaled.png"
        alt="Saferwall"
        class="logo-horizontal"
      />
      <img
        src="../../assets/imgs/saferwall_jkdgq6_c_scale_w_800.png"
        alt="Saferwall"
        class="logo-square"
      />
    </router-link>
    <div class="mobile-nav" @click="showinmobile = !showinmobile">
      <i class="icon ion-android-menu"></i>
    </div>
    <div class="header-search" :class="{ active: showinmobile }">
      <input
        type="search"
        placeholder="Quick file hash (sha256) lookup"
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
          <router-link :to="this.$routes.UPLOAD.path" class="button-upload">
            <i class="icon fas fa-upload fa-2x"></i>
          </router-link>
        </li>
        <li>
          <a href="https://about.saferwall.com/">
            <span>
              About
            </span>
          </a>
        </li>
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
                  <span>{{ getLoggedIn ? getUsername : "Login" }}</span>
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
                    <router-link :to="this.$routes.PROFILE.path+getUsername">
                      <i class="icon fas fa-user-circle fa-lg"></i>
                      Profile
                    </router-link>
                  </div>
                  <div class="dropdown-item">
                    <router-link :to="this.$routes.SETTINGS.path">
                      <i class="icon fas fa-cogs "></i>
                      Settings
                    </router-link>
                  </div>
                  <div class="dropdown-item">
                    <div @click="loginOrLogout" class="has-text-danger">
                      <span>
                        <i class="icon fas fa-sign-out-alt fa-lg"></i
                        > Logout</span
                      >
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
    ...mapActions(["logOut", "updateHash"]),
    showMobileSearch() {},
    loginOrLogout(e) {
      if (this.getLoggedIn && e.target.textContent.trim() === "Logout") {
        this.logOut()
        this.$router.go(this.$routes.HOME.path)
      } else if (!this.getLoggedIn) {
        this.$router.push(this.$routes.LOGIN.path)
      }
    },
    searchByHash() {
      if (!this.hash.trim()) {
        this.$awn.warning("Empty Field!")
        return
      }
      this.$http
        .get(`${this.$api_endpoints.FILES}${this.hash}/`, {
          validateStatus: (status) => status === 200,
        })
        .then((data) => {
          this.updateHash(this.hash)
          this.track()
          this.$router.push(this.$routes.SUMMARY.path + this.hash)
        })
        .catch(() => {
          this.$awn.alert(
            "Sorry, we couldn't find the file you were looking for, please upload it to view the results!",
          )
        })
    },
    track() {
      this.$gtag.event("search", {
        search_term: this.hash,
      })
    },
  },
}
</script>

<style scoped lang="scss">
@import "../../assets/scss/variables";
@import url("https://cdn.jsdelivr.net/npm/@mdi/font@4.x/css/materialdesignicons.min.css");
header.dashboard-header {
  margin-right: 20px;
  background: #fff;
  height: $header-height;
  line-height: $header-height;
  // box-shadow: 0 2px 3px rgba(10, 10, 10, 0.1), 0 0 0 1px rgba(10, 10, 10, 0.1);
  border-bottom: 1px solid rgba(25,25,25,0.1);
  z-index: 102;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  display: flex;
  .logo-square{
    display: none;
  }
  .logo {
    float: left;
    line-height: 62px;
    height: 80px;
    width: 250px;
    text-align: center;
    margin: auto;
    display: grid;
    padding: 0 20px;

    img {
      width: 100%;
      padding: 15px;
      height: auto;
    }
  }

  .mobile-nav {
    // margin-left:20px;
    float: right;
    height: $header-height;
    width: $header-height;
    line-height: $header-height;
    border-left: solid 1px rgba(10, 10, 10, 0.1);
    border-right: solid 1px rgba(10, 10, 10, 0.1);
    text-align: center;
    font-size: 2rem;

    @media screen and (min-width: 792px) {
      display: none;
    }
  }

  .header-search {
    display: inline-block;
    width: 500px;
    height: $header-height;
    position: relative;
    padding-left: 15px;
    border-left: solid 1px rgba(black, 0.1);
    flex: 1;
    margin-right: 30px;

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
  .header-search input{
    border-bottom: 1px solid rgba(25,25,25,0.1);
    outline: none !important;
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

        *{
          max-height: $header-height !important;
        }
        // border-left: solid 1px rgba(black, 0.1);
        // min-width: 90px;
        text-align: center;
        a {
          display: inline-block;
          margin: 0 10px;
          font-size: 14px;
          color: #2c3e50;
          font-weight: 400;
          transition: all 0.2s;

          &:hover {
            color: $primary-color;
          }
        }
        span{
          padding: 10px 20px;
          background: #f7f7f7;
          border-radius: 3px;
        }
        .button-upload{
          border-radius: 3px;
          padding: 0 50px;
          background: #f7f7f7;
          border-bottom: rgba(25,25,25,0.1) solid 1px;
        }
        & > .profile {
          cursor: pointer;
          padding-left: 10px;
          padding-right: 10px;
          font-weight: 500;

          .has-text-danger {
            display: inline-block;
            font-size: small;

            &:hover {
              font-weight: 600;
            }
          }

          .dropdown-content {
            text-align: left;
            a {
              margin: 0;
            }
          }

          .dropdown-menu {
            min-width: 10rem;
          }

          .dropdown_text,
          .has-text-link {
            color: #2c3e50 !important;
            font-weight: 400 !important;
            font-size: 14px;

            & > span:hover {
              color: $primary-color;
            }
          }
        }

        button {
          display: inline-block;
          width: 100%;
          color: rgba(black, 0.6);
          font-weight: 600;
          font-size: 14px;
          text-align: left;
          padding: 0 5px 0 5px;
          transition: all 0.2s;
          background: none;
          border: none;
          cursor: pointer;
          &:hover {
            color: $primary-color !important;
          }
        }

        .fa-upload {
          color: $primary-color !important;
          &:hover {
            color: #2c3e50 !important;
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
      z-index: 102;
      border-bottom: 1px solid rgba(25,25,25,0.1);
      // box-shadow: 0 5px 10px rgba(10, 10, 10, 0.1);
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
.fa-cogs{
  vertical-align: bottom;
}
@media (min-width: 795px) {
  .minimized .dashboard-header {
    .logo{
      width: $header-height - 1px; 
      padding: 0;
    }
    .logo-horizontal{
      display: none;
    }
    .logo-square{
      display: block;
    }
    img{
      padding :0 !important;
    }
  }
}
</style>
