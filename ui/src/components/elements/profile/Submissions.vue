<template>
  <div v-if="active">
    <p id="no" v-if="!subs">No Submissions</p>
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
    getFileData: function(file) {
      this.$http
        .get(
          this.$api_endpoints.FILES +
            file.sha256 +
            "?fields=sha256,tags,multiav",
        )
        .then((res) => {
          res.data.AvDetectionCount = this.getAvDetectionCount(
            res.data.multiav.last_scan,
          )
          res.data.timestamp = file.timestamp
          this.filesData.push(res.data)
        })
        .catch()
    },
    loadFiles: function() {
      for (var index in this.subs) {
        this.getFileData(this.subs[index])
      }
    },
  },
  mounted() {
    this.loadFiles()
  },
}
</script>

<style lang="scss" scoped>
#no {
  font-size: 25px;
  font-weight: 200;
  padding: 0.5em;
}</style>
