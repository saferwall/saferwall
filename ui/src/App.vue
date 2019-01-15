<template>
  <div id="app">
    <v-layout>
      <router-view/>
    </v-layout>
  </div>
</template>

<script>
import Layout from '@/components/Layout'
import router from 'vue-router'
export default {
  name: 'App',
  components: {
    'v-layout': Layout,
  },
  created(){
    document.title = this.$route.meta.title
  },
  updated(){
    document.title = this.$route.meta.title
  }
}
</script>

<style lang="scss">
@import '../node_modules/bulma/bulma';
@import 'assets/scss/variables';
@import 'assets/scss/typography';
@import 'assets/scss/layout';
@import 'assets/scss/ionicons';

*{
  outline:none!important;
  padding:0;
  margin:0;
}

html{
  padding:0;
  margin:0!important;
  background-color: rgba(0,0,0,.03);
  font-size:14px;
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
</style>
