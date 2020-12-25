<template>
  <section class="main-content" :class="{ fullwidth: fullwidth }">
    <div class="container is-fluid">
      <div class="columns top_columns">
        <div class="column is-9 box page-path">
          <nav class="breadcrumb" aria-label="breadcrumbs">
            <ul>
              <li>
                <router-link :to="this.$routes.HOME.path">Home</router-link>
              </li>
              <li class="is-active" v-if="route !== 'home'">
                <a href="#" aria-current="page">{{ route }}</a>
              </li>
              <li v-if="$route.params.hash">
                <span>{{ $route.params.hash }}</span>
              </li>
            </ul>
          </nav>
        </div>
        <div class="column no-shadow" v-if="showButtons">
          <div class="buttons" >
            <Like :hash="hash" />
            <Download :hash="hash" />
            <Rescan :route="route" :hash="hash" />
          </div>
        </div>
      </div>
      <div v-if="!showContent" class="column ">
        <p class="no_file placeholders" >No file Specified</p>
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
  padding-top: $content-marging;
  margin-top: $header-height;
  margin-left: $sidebar-width + $content-marging;
  
  &:not(.fullwidth) {
    width: calc(100% - 60px * 2 );
  }

  &.fullwidth {
    width: 100%;
  }
}
.container{
  padding: 0;
}
.breadcrumb {
  padding: 10px 20px;
  background: #fff;
  border-radius: 3px;

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
.buttons {
  margin-left: 5px;
  margin-bottom : unset;
  justify-content: center;
  display: flex;
  align-items: center;
  width: 100%;
  
  *{
    min-width: calc(100%/3.2);
    margin: auto;
  }
  .button{
    min-height: 41px;
    padding: 0;
  }
}

.top_columns{
  margin-bottom: 20px;
  max-width: 100%;

  .column{
    margin: 0;
    padding: 0;
    display:flex;
  }
  .page-path{
    box-shadow: rgba(25,25,25,0.1) 1px 1px 5px;
  }
}
</style>
