<template>
  <div>
    <v-header></v-header>
    <transition name="slideLeft">
      <v-sidebar v-if="noSidebarPaths.indexOf($route.name) == -1"></v-sidebar>
    </transition>
    <v-content :fullwidth="noSidebarPaths.indexOf($route.name) > -1">
      <router-view></router-view>
    </v-content>
    <v-footer v-if="false"></v-footer>
  </div>
</template>
<script>
import Header from "@/components/partials/Header"
import Footer from "@/components/partials/Footer"
import Sidebar from "@/components/partials/Sidebar"
import Content from "@/components/partials/Content"

export default {
  components: {
    "v-header": Header,
    "v-sidebar": Sidebar,
    "v-content": Content,
    "v-footer": Footer,
  },

  data() {
    return {
      noSidebarPaths: ["Upload", "Scanning", "settings", "profile", "home"],
    }
  },
}
</script>
<style lang="scss" scoped>
.fade-enter-active,
.fade-leave-active {
  transition-property: opacity;
  transition-duration: 0.25s;
}

.fade-enter-active {
  transition-delay: 0;
}

.fade-enter,
.fade-leave-active {
  opacity: 0;
}

.slideLeft-enter-active,
.slideLeft-leave-active {
  transition-property: transform;
  transition-duration: 0.25s;
}

.slideLeft-enter,
.slideLeft-leave-active {
  transform: translateX(-100%);
}
</style>
