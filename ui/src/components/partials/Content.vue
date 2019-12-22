<template>
  <section class="main-content" :class="{ fullwidth: fullwidth }">
    <transition name="slide" mode="out-in">
      <notification :type="notifType" @closeNotif="close()" v-if="notifActive">
        {{ notifText }}
      </notification>
    </transition>
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
          <Download v-if="showDownload" />
        </div>
        <div class="column">
          <Rescan v-if="showRescan" :route="route" @notify="showNotification" />
        </div>
      </div>
      <slot @notify="showNotification" v-if="showContent"></slot>
    </div>
  </section>
</template>
<script>
import Download from "../elements/Download"
import Rescan from "../elements/Rescan"
import Notification from "../elements/Notification"

export default {
  props: ["fullwidth"],
  data() {
    return {
      route: "",
      notifText: "",
      notifActive: false,
      notifType: "",
    }
  },
  components: {
    Download,
    Rescan,
    Notification,
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
    showContent : function(){
      return this.$store.getters.getHashContext || this.route === "upload"
    }
  },
  methods: {
    close: function() {
      this.notifActive = false
    },
    showNotification: function(type, text) {
      this.notifActive = true
      this.notifType = type
      this.notifText = text
    },
  },
  created() {
    this.route = this.$router.currentRoute.name
    if (
      this.$router.currentRoute.params.hash &&
      this.$router.currentRoute.params.hash !==
        this.$store.getters.getHashContext
    ) {
      this.$http
        .get(this.$api_endpoints.FILES + this.$router.currentRoute.params.hash)
        .then((data) => {
          this.$store.dispatch(
            "updateHash",
            this.$router.currentRoute.params.hash,
            this.$store.dispatch("updateFileData", data),
          )
        })
    } else if (!this.$router.currentRoute.params.hash)
      this.$store.dispatch("updateHash", "")
  },
  updated() {
    this.route = this.$router.currentRoute.name
    if (
      this.$router.currentRoute.params.hash &&
      this.$router.currentRoute.params.hash !==
        this.$store.getters.getHashContext
    ) {
      this.$http
        .get(this.$api_endpoints.FILES + this.$router.currentRoute.params.hash)
        .then((data) => {
          this.$store.dispatch(
            "updateHash",
            this.$router.currentRoute.params.hash,
          )
          this.$store.dispatch("updateFileData", data)
        })
    } else if (
      !this.$router.currentRoute.params.hash &&
      !this.$store.getters.getHashContext
    )
      this.$store.dispatch("updateHash", "")
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
</style>
