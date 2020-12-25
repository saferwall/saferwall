<template>
  <section class="main-content" :class="{ fullwidth: fullwidth }">
    <div class="container is-fluid">
      <div class="columns top_columns">
        <div class="column is-9">
          <nav class="breadcrumb" aria-label="breadcrumbs">
            <ul>
              <li>
                <router-link :to="this.$routes.HOME.path">Home</router-link>
              </li>
              <li class="is-active" v-if="route !== 'home'">
                <a href="#" aria-current="page">{{ route }}</a>
              </li>
              <li>
                <span>{{ $route.params.hash }}</span>
              </li>
            </ul>
          </nav>
        </div>
        <div class="column ">
          <div class="buttons" v-if="showButtons">
            <Like :hash="hash" />
            <Download :hash="hash" />
            <Rescan :route="route" :hash="hash" />
          </div>
        </div>
      </div>
      <div class="column placeholders">
        <p class="no_file" v-if="!showContent">No file Specified</p>
        <loader v-if="false"></loader>
      </div>
      <slot v-if="showContent"></slot>
      <div class="column">
        <Social />
      </div>
    </div>
  </section>
</template>

<script>
import Loader from "@/components/elements/Loader"
import Download from "../elements/Download"
import Rescan from "../elements/Rescan"
import Social from "../elements/Social"
import Like from "../elements/LikeButton"

import { mapGetters } from "vuex"

export default {
  props: ["fullwidth"],
  data() {
    return {
      route: "",
    }
  },
  components: {
    Download,
    Rescan,
    Social,
    Like,
    Loader,
  },
  computed: {
    ...mapGetters({
      hash: "getHashContext",
      userData: "getUserData",
      fileData: "getFileData",
      username: "getUsername",
      loggedIn: "getLoggedIn",
    }),
    showButtons: function() {
      return (
        Object.entries(this.$store.getters.getFileData).length !== 0 &&
        this.route !== "upload" &&
        this.route !== "settings" &&
        this.route !== "home" &&
        this.route !== "profile"
      )
    },
    showContent: function() {
      if (!this.loggedIn) {
        return (
          this.$store.getters.getHashContext !== "" ||
          this.route === "upload" ||
          this.route === "profile" ||
          this.route === "home" ||
          this.route === "settings"
        )
      } else {
        return (
          this.$store.getters.getHashContext !== "" ||
          this.route === "upload" ||
          this.route === "profile" ||
          this.route === "home" && this.username||
          this.route === "settings"
        )
      }
    },
  },
  methods: {
    getData: function() {
      this.route = this.$router.currentRoute.name
      if (
        this.$router.currentRoute.params.hash &&
        this.$router.currentRoute.params.hash !==
          this.$store.getters.getHashContext
      ) {
        this.$store.dispatch(
          "updateHash",
          this.$router.currentRoute.params.hash,
        )
      }
    },
  },
  created() {
    this.getData()
  },
  updated() {
    this.route = this.$router.currentRoute.name
  },
}
</script>

<style scoped lang="scss">
@import "../../assets/scss/variables";

section.main-content {
  float: unset;
  padding-top: 20px;
  margin-top: $header-height;
  margin-left: $sidebar-width + 20px;

  &:not(.fullwidth) {
    width: calc(100% - 200px);
  }

  &.fullwidth {
    width: 100%;
  }
}
.breadcrumb {
  a {
    color: $primary-color;
  }
  span {
    padding-left: 10px;
  }
}
.no_file {
  font-size: 20px;
  font-weight: 200;
}
#loader {
  margin-top: 1em;
  margin-bottom: 2em;
}
.placeholders {
  margin-bottom: 2em;
}
.buttons *{
  display: block;
  text-align:right;
}
</style>
