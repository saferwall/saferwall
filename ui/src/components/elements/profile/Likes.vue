<template>
  <div class="likes">
    <div class="tile like is-child box" v-for="(file, index) in filesData" :key="index">
        <div id="hash">{{file.sha256}}</div>
    </div>
  </div>
</template>

<script>
export default {
  props: ["likes"],
  data() {
    return {
      filesData: [],
    }
  },
  methods: {
    getFileData: function(hash) {
      this.$http
        .get(this.$api_endpoints.FILES + hash+"?fields=sha256,tags,multiav")
        .then((res) => {
          this.filesData.push(res.data)
        })
        .catch()
    },
  },
  mounted() {
    for (var index in this.likes) {
      this.getFileData(this.likes[index])
    }
  },
}
</script>

<style></style>
