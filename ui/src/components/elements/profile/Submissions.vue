<template>
  <div class="likes" v-if="active">
    <FileCard :file="file" v-for="(file, index) in filesData" :key="index" />
  </div>
</template>

<script>
import FileCard from "./FileCard"

export default {
  props: ["subs", "active"],
  components: {
    FileCard,
  },
  data() {
    return {
      filesData: [],
    }
  },
  watch: {
    subs: function() {
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
      for (var index in this.subs) {
        this.getFileData(this.subs[index].hash)
      }
    },
  },
  mounted() {
    this.loadFiles()
  },
}
</script>

<style lang="scss" scoped></style>
