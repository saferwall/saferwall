<template>
  <div>
    <loader v-if="showLoader"></loader>
    <div class="tile is-ancestor" v-if="!showLoader">
      <div class="tile is-parent is-vertical">
        <!-- Basic Properties -->
        <basicProperties :summaryData="summaryData" />
        <!--ExifTool File Metadata-->
        <div class="tile is-parent">
          <exifTool :summaryData="summaryData" />
          <submissions :summaryData="summaryData" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Loader from "@/components/elements/Loader"

import BasicProperties from "../elements/summary/BasicProperties"
import ExifTool from "../elements/summary/ExifTool"
import Submissions from "../elements/summary/Submissions"

import { mapGetters } from "vuex"

export default {
  components: {
    loader: Loader,
    basicProperties: BasicProperties,
    exifTool: ExifTool,
    submissions: Submissions,
  },
  computed: {
    ...mapGetters({ fileData: "getFileData" }),
    summaryData: function() {
      if (
        !this.fileData ||
        (Object.entries(this.fileData).length === 0 &&
          this.fileData.constructor === Object)
      )
        return {}
        
      return this._.omit(this.fileData.data, ["multiav", "strings", "status"])
    },
    showLoader: function() {
      return this.summaryData === {} || !this.summaryData
    },
  },
}
</script>

<style lang="scss">
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
    &.exif {
      width: 30%;
    }
    &.subs {
      width: 20%;
    }
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
.tile.is-parent{
  padding: 0;
  margin-bottom: 10px;
}
.tile.is-child {
    margin-right: 18px !important;
}
</style>
