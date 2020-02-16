<template>
  <div class="likes" v-if="active">
    <LikedFile :file="file" v-for="(file, index) in filesData" :key="index" />
  </div>
</template>

<script>
import LikedFile from "./LikedFile"

export default {
  props: ["likes", "active"],
  components: {
    LikedFile,
  },
  data() {
    return {
      filesData: [],
    }
  },
  watch: {
    likes: function() {
      this.filesData = []
      this.loadFiles()
    },
  },
  methods: {
    getAvDetectionCount: function(scans) {
      var count = 0
      for (const av of Object.values(scans)) {
        if (av.infected) count++
      }
      return count
    },
    getFileData: function(hash) {
      this.$http
        .get(this.$api_endpoints.FILES + hash + "?fields=sha256,tags,multiav")
        .then((res) => {
          res.data.AvDetectionCount = this.getAvDetectionCount(
            res.data.multiav.last_scan,
          )
          this.filesData.push(res.data)
        })
        .catch()
    },
    loadFiles: function() {
      for (var index in this.likes) {
        this.getFileData(this.likes[index])
      }
    },
  },
  mounted() {
    this.loadFiles()
  },
}
</script>

<style lang="scss" scoped></style>
