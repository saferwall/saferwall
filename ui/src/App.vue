<template>
  <div id="app">
    <transition name="component-fade" mode="out-in">
    <component :is="layout">
      <router-view />
    </component>
    </transition>
  </div>
</template>

<script>
import router from "vue-router";
export default {
  name: "App",
  computed: {
    layout() {
      return (this.$route.meta.layout || "default") + "-layout";
    }
  },
  created() {
    document.title = this.$route.meta.title;
  },
  updated() {
    document.title = this.$route.meta.title;
  }
};
</script>

<style lang="scss">
@import "../node_modules/bulma/bulma";
@import "assets/scss/variables";
@import "assets/scss/typography";
@import "assets/scss/layout";
@import "assets/scss/ionicons";

* {
  padding: 0;
  margin: 0;
}

*::selection {
  background-color: #18a096;
  color: white;
}

html {
  padding: 0;
  margin: 0 !important;
  background-color: rgba(0, 0, 0, 0.03);
  font-size: 14px;
  font-weight: 400;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  // font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

.component-fade-enter-active, .component-fade-leave-active {
  transition: opacity .2s ease;
}
.component-fade-enter, .component-fade-leave-to
 {
  opacity: 0;
}
</style>
