<template>
  <section class="main-content" :class="{ fullwidth: fullwidth }">
    <div class="container is-fluid">
      <div class="columns">
        <div
          class="column is-four-fifths-fullhd is-three-quarters-desktop is-three-fifths-tablet is-half-mobile"
        >
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
        <div class="column">
            <div class="buttons">
              <Download v-if="showDownload" :hash="hash"/>
              <Rescan v-if="showRescan" :route="route" :hash="hash" />
            </div>
        </div>
      </div>
      <p class="no_file" v-if="!showContent">No file Specified</p>
      <slot v-if="showContent"></slot>
    </div>
  </section>
</template>
<script>
import Download from "../elements/Download"
import Rescan from "../elements/Rescan"

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
  },
  computed: {
    showDownload: function() {
      return (
        this.$store.getters.getHashContext &&
        this.$store.getters.getLoggedIn &&
        this.route !== "upload"
      )
    },
    showRescan: function() {
      return (
        this.$store.getters.getHashContext && this.$store.getters.getLoggedIn
      )
    },
    showContent: function() {
      return (
        this.$store.getters.getHashContext !== "" || this.route === "upload"
      )
    },
    ...mapGetters({
      hash: 'getHashContext'
    }),
  },
  methods: {
    getData: function() {
      this.route = this.$router.currentRoute.name
      if (
        this.$router.currentRoute.params.hash &&
        this.$router.currentRoute.params.hash !==
          this.$store.getters.getHashContext
      ) {
        this.$http
          .get(
            this.$api_endpoints.FILES + this.$router.currentRoute.params.hash,
          )
          .then((data) => {
            this.$store.dispatch(
              "updateHash",
              this.$router.currentRoute.params.hash,
            )
            this.$store.dispatch("updateFileData", data)
          })
          .catch(() => {
            this.$awn.alert(
              "Sorry, we couldn't find the file you were looking for, please upload it to view the results!",
            )
          })
      }
    },
  },
  created() {
    this.getData()
  },
  updated() {
    this.getData()
  },
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
</style>
