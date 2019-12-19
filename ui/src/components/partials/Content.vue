<template>
  <section class="main-content" :class="{ fullwidth: fullwidth }">
    <div class="container is-fluid">
      <div class="columns">
        <div class="column is-four-fifths">
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
          <Download v-if="show"/>
        </div>
        <div class="column">
          <Rescan v-if="show"/>
        </div>
      </div>
      <slot></slot>
    </div>
  </section>
</template>
<script>
import Download from "../elements/Download"
import Rescan from "../elements/Rescan"
export default {
  props: ["fullwidth"],
  data() {
    return {
      route : ""
    }
  },
  components: {
    Download,
    Rescan,
  },
  computed : {
    show : function(){
      return (this.$store.getters.getHashContext && this.$store.getters.getLoggedIn && this.route !== "upload")
    }
  },
  created() {
    this.route = this.$router.currentRoute.name
    if (this.$router.currentRoute.params.hash)
      this.$store.dispatch("updateHash", this.$router.currentRoute.params.hash)
    else this.$store.dispatch("updateHash", "")
  },
  updated(){
    this.route = this.$router.currentRoute.name
  }
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
</style>
