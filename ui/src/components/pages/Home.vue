<template>
    <p>
        <router-link to="/upload">go to upload</router-link>
    </p>
</template>

<script>
export default {
  name: 'Home',
  data () {
    return {
    }
  },

  mounted() {
    console.log(process.env.API_ROOT)
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
