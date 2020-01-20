<template>
  <div>
    <loader v-if="showLoader"></loader>
    <div class="tile is-ancestor" v-if="!showLoader">
      <div class="tile is-parent is-vertical">
        <div class="tile is-child box">
          <h4 class="title">Basic Properties</h4>
          <!--Basic Properties-->
          <div
            v-for="(i, index) in basicProperties"
            class="data-data"
            :key="index"
          >
            <strong class="data-label">
              {{ getLabelForGivenKey(index) }}
            </strong>
            <!-- TRiD -->
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
            <!-- Packer -->
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
            <!-- Tags -->
            <span class="tags" v-else-if="index === 'tags'">
              <span
                v-for="(item, index) in summaryData.tags"
                :key="index"
                class="tag is-link is-normal"
              >
                <span>{{ item }}</span>
                <copy :content="item"></copy>
              </span>
            </span>
            <!-- Default -->
            <span class="data-value" v-else>
              <span class="value-text">{{
                index !== "sha512" ? i : i.substring(0, 70) + "..."
              }}</span>
              <copy :content="i"></copy>
            </span>
          </div>
        </div>

        <!--ExifTool File Metadata-->
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

import { mapGetters } from "vuex"

export default {
  components: {
    loader: Loader,
    copy: Copy,
  },
  data() {
    return {
      uppercaseFields: ["md5", "sha-1", "sha-256", "sha-512", "crc32"],
    }
  },
  computed: {
    ...mapGetters({ fileData: "getFileData" }),
    basicProperties: function() {
      const allPropsEntries = Object.entries(this.summaryData)
      const basicPropsEntries = allPropsEntries.filter(
        (entry) => !["av", "exif"].includes(entry[0]),
      )
      return Object.fromEntries(basicPropsEntries)
    },
    summaryData: function() {
      if (!this.fileData || Object.entries(this.fileData).length === 0 && this.fileData.constructor === Object) return {}
      return {
        filesize: this.bytesToSize(this.fileData.data.size),
        magic: this.fileData.data.magic,
        crc32: this.fileData.data.crc32,
        md5: this.fileData.data.md5,
        "sha-1": this.fileData.data.sha1,
        "sha-256": this.fileData.data.sha256,
        "sha-512": this.fileData.data.sha512,
        ssdeep: this.fileData.data.ssdeep,
        trid: this.fileData.data.trid,
        packer: this.fileData.data.packer,
        tags: this.fileData.data.tags,
        exif: this.fileData.data.exif,
      }
    },
    showLoader: function() {
      return this.summaryData === {} || !this.summaryData
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

  .tag {
    position: relative;

    & > .copy {
      opacity: 0;
    }

    &:hover {
      span {
        opacity: 0.35;
      }
      & > .copy {
        opacity: 1;
      }
    }
  }
}
</style>
