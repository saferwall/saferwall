<template>
  <section class="main-content" :class="{ fullwidth: fullwidth }">
    <div class="container is-fluid">
      <div class="columns top_columns">
        <div class="column is-9">
          <nav class="breadcrumb" aria-label="breadcrumbs" v-if="!fullwidth">
            <ul>
              <li>
                <router-link :to="this.$routes.HOME.path"
                  >Dashboard</router-link
                >
              </li>
              <li class="is-active">
                <a href="#" aria-current="page">{{ $route.name }}</a>
              </li>
            </ul>
          </nav>
        </div>
        <div class="column is-4">
          <div class="buttons" v-if="!showLoader && showButtons">
            <Like :hash="hash" />
            <Download :hash="hash" />
            <Rescan :route="route" :hash="hash" />
          </div>
        </div>
      </div>
      <div class="column placeholders">
        <p class="no_file" v-if="!showContent">No file Specified</p>
        <loader v-if="showLoader && showContent"></loader>
      </div>
      <slot v-if="showContent && !showLoader"></slot>
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
    }),
    showButtons: function() {
      return (
        Object.entries(this.$store.getters.getFileData).length !== 0 &&
        this.$store.getters.getLoggedIn &&
        this.route !== "upload" &&
        this.route !== "profile"
      )
    },
    showContent: function() {
      return (
        this.$store.getters.getHashContext !== "" ||
        this.route === "upload" ||
        this.route === "profile"
      )
    },
    showLoader: function() {
      if (this.$store.getters.getLoggedIn)
        return (
          Object.entries(this.userData).length === 0 &&
          this.userData.constructor === Object
        )
      else
        return (
          Object.entries(this.fileData).length === 0 &&
          this.fileData.constructor === Object
        )
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
  // updated() {
  //   this.getData()
  // },
}
</script>
<style scoped lang="scss">
@import "../../assets/scss/variables";
$header-height: 50px;
section.main-content {
  float: right;
  padding-top: 20px;
  margin-top: $header-height;

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
</style>
