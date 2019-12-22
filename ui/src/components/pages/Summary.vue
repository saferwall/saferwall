<template>
  <div>
    <loader v-if="showLoader"></loader>
    <div class="tile is-ancestor" v-if="!showLoader">
      <div class="tile is-parent is-vertical">
        <div class="tile is-child box">
          <h4 class="title">Basic Properties</h4>
          <div
            v-for="(i, index) in basicProperties"
            class="data-data"
            :key="index"
          >
            <strong class="data-label">
              {{ getLabelForGivenKey(index) }}
            </strong>
            <span
              class="data-value"
              :class="{ 'trid-container': index === 'trid' }"
              v-if="index === 'trid'"
            >
              <p v-for="(t, index) in summaryData.trid" :key="index">
                <span class="trid">
                  <span class="value-text">{{ t }}</span>

                  <copy :content="t"></copy>
                </span>
              </p>
            </span>
            <span
              class="data-value"
              :class="{ 'packer-container': index === 'packer' }"
              v-else-if="index === 'packer'"
            >
              <p v-for="(t, index) in summaryData.packer" :key="index">
                <span class="packer">
                  <span class="value-text">{{ t }}</span>

                  <copy :content="t"></copy>
                </span>
              </p>
            </span>
            <span class="data-value" v-else>
              <span class="value-text">{{
                index !== "sha512" ? i : i.substring(0, 70) + "..."
              }}</span>

              <copy :content="i"></copy>
            </span>
          </div>
        </div>

        <div class="tile is-child box">
          <h4 class="title">ExifTool File Metadata</h4>
          <div
            v-for="(i, index) in summaryData.exif"
            :key="index"
            class="data-data"
          >
            <strong class="data-label">
              {{ index.replace(/[A-Z]/g, (match) => ` ${match}`) }}
            </strong>
            <span class="data-value">
              <span class="value-text">{{ i }}</span>

              <copy :content="i"></copy>
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import Loader from "@/components/elements/Loader"
import Copy from "@/components/elements/Copy"

export default {
  components: {
    loader: Loader,
    copy: Copy,
  },
  data() {
    return {
      showLoader: true,
      summaryData: {},
      uppercaseFields: ["md5", "sha-1", "sha-256", "sha-512", "crc32"],
    }
  },
  computed: {
    basicProperties: function() {
      const allPropsEntries = Object.entries(this.summaryData)
      const basicPropsEntries = allPropsEntries.filter(
        (entry) => !["av", "exif"].includes(entry[0]),
      )
      return Object.fromEntries(basicPropsEntries)
    },
  },
  methods: {
    bytesToSize(bytes) {
      var sizes = ["Bytes", "KB", "MB", "GB", "TB"]
      if (bytes === 0) return "0 Byte"
      var i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)))
      return Math.round(bytes / Math.pow(1024, i), 2) + " " + sizes[i]
    },
    getLabelForGivenKey(key) {
      return this.uppercaseFields.includes(key)
        ? key.toUpperCase()
        : key === "filesize"
        ? "File Size"
        : key === "trid"
        ? "TRiD"
        : key === "ssdeep"
        ? "SSDeep"
        : key
    },
    showData() {
      var fileData = this.$store.getters.getFileData
      if(fileData === {} || !fileData) return
      this.showLoader = false

      fileData.data["sha-1"] = fileData.data.sha1
      fileData.data["sha-256"] = fileData.data.sha256
      fileData.data["sha-512"] = fileData.data.sha512
      delete fileData.data.sha1
      delete fileData.data.sha256
      delete fileData.data.sha512

      this.summaryData.filesize = this.bytesToSize(fileData.data.size)
      this.summaryData.magic = fileData.data.magic
      this.summaryData.crc32 = fileData.data.crc32
      this.summaryData.md5 = fileData.data.md5
      this.summaryData["sha-1"] = fileData.data["sha-1"]
      this.$set(this.summaryData, "sha-256", fileData.data["sha-256"]) // Setting a reactive property (Object entries are not reactive), this is used to update BasicProperties
      this.summaryData["sha-512"] = fileData.data["sha-512"]
      this.summaryData.ssdeep = fileData.data.ssdeep
      this.summaryData.trid = fileData.data.trid
      this.summaryData.exif = fileData.data.exif
      this.summaryData.packer = fileData.data.packer
    },
  },
  mounted() {
    if (this.$store.getters.getHashContext) this.showData()
  },
  beforeRouteUpdate(to, from, next) {
    this.showData()
  },
}
</script>
<style lang="scss" scoped>
.data-data {
  float: left;
  width: 100%;
  padding: 5px;

  &:nth-child(even) {
    background: rgba(black, 0.03);
  }

  &:not(:last-child),
  .trid:not(:last-child) {
    margin-bottom: 3px;
  }

  .data-label {
    float: left;
    width: 10.5%;
    text-transform: capitalize;
    margin-right: 1.4em;
  }

  .data-value {
    float: left;

    .value-text {
      transition: opacity 0.2s;
    }

    .copy {
      opacity: 0;
      transition: opacity 0.2s;
    }

    &:not(.trid-container):not(.packer-container):hover {
      .value-text {
        opacity: 0.35;
      }
      & > .copy {
        opacity: 1;
      }
    }
  }

  .trid,
  .packer,
  .data-value {
    position: relative;
  }

  .trid,
  .packer {
    position: relative;

    &:hover {
      .value-text {
        opacity: 0.35;
      }
      & > .copy {
        opacity: 1;
      }
    }
  }
}
</style>
