<template>
  <div class="container" v-if="this.data">
    <div class="section_title">Functions</div>
    <div class="sections_header">
      <div class="sections_header_field" v-for="label in labels" :key="label">
        {{ _.startCase(label) }}
      </div>
    </div>
    <div
      class="section"
      v-for="(section, index) in filteredData"
      :key="index + 'b'"
    >
      <div class="section_content">
        <div
          class="section_field"
          v-for="(field, index) in Object.values(section)"
          :key="'function_' + index"
        >
          <span class="parent">
            <span>
              {{ toHex(field) }}
            </span>
            <copy :content="toHex(field)" />
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { dec2HexString } from "@/helpers/pe"
import Copy from "@/components/elements/Copy"

export default {
  props: ["data"],
  components: {
    copy: Copy,
  },
  computed: {
    byOrdinal: function() {
      if (this.data) return this.data[0]["ByOrdinal"]
      else return false
    },
    labels: function() {
      if (this.data) {
        var labelsTmp = Object.keys(this.data[0])
        if (!this.byOrdinal) return this._.without(labelsTmp, "ByOrdinal", "Ordinal")
        else return this._.without(labelsTmp, "ByOrdinal")
      }
      return []
    },
    filteredData: function() {
      if (!this.byOrdinal) {
        return this._.map(this.data, (obj) => {
          return this._.omit(obj, ["ByOrdinal", "Ordinal"])
        })
      }
      return this._.map(this.data, (obj) => {
        return this._.omit(obj, ["ByOrdinal"])
      })
    },
  },
  methods: {
    toHex: function(value) {
      if (this._.isNumber(value)) return dec2HexString(value)
      return value
    },
  },
}
</script>

<style lang="scss" scoped>
.container {
  overflow: hidden;
  margin-top: 2rem;
  .section_title {
    font-size: large;
    color: #00d1b2;
  }
}
.sections_header {
  display: inline-flex;
  padding: 0.2rem;
  .sections_header_field {
    text-align: left;
    width: 17rem;
    font-weight: 600;
    &.name {
      width: 6rem;
    }
  }
}
.section {
  padding: 0.2rem;
  .section_content {
    display: inline-flex;
    .section_field {
      text-align: left;
      width: 17rem;
      &.name {
        width: 6rem;
      }
      &:hover {
        .copy {
          opacity: 1;
        }
      }
    }
  }
}
.parent {
  position: relative;
  .copy {
    opacity: 0;
    transition: opacity 0.2s;
  }
}
</style>
