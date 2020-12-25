<template>
  <div>
    <div class="columns is-1">
      <div class="column is-6">
        <basicProperties :summaryData="summaryData" />
      </div>
      <div class="column is-6">
        <exifTool :summaryData="summaryData" />
      </div>
    </div>
    <div class="columns">
      <div class="column">
        <submissions :summaryData="summaryData" />
      </div>
    </div>
  </div>
</template>

<script>
import BasicProperties from "../elements/summary/BasicProperties"
import ExifTool from "../elements/summary/ExifTool"
import Submissions from "../elements/summary/Submissions"

import { mapGetters } from "vuex"

export default {
  components: {
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

      return this._.omit(this.fileData, ["multiav", "strings", "status"])
    },
  },
}
</script>

<style lang="scss">
.data-data {
  width: 100%;
  padding: 5px;
  display: flex;

  &:nth-child(even) {
    background: rgba(black, 0.03);
  }

  &:not(:last-child),
  .trid:not(:last-child) {
    margin-bottom: 3px;
  }

  .data-label {
    min-width: 100px;
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

    &.redTag {
      background-color: #f14668;
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
.tile.is-parent {
  padding: 0;
  margin-bottom: 10px;
}
</style>
