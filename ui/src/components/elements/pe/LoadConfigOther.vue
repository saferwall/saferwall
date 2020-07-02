<template>
  <div class="container">
    <div class="sections_header">
      <div class="sections_header_field" v-for="label in labels" :key="label">
        {{ _.startCase(label) }}
      </div>
    </div>
    <div
      class="section"
      v-for="(section, index) in selectedSection"
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
              {{ processData(field) }}
            </span>
            <copy :content="processData(field)" />
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Copy from "@/components/elements/Copy"
import { dec2HexString } from "@/helpers/pe"

export default {
  props: ["data", "section"],
  components: {
    copy: Copy,
  },
  computed: {
    sections: function() {
      return this._.omitBy(Object.keys(this.data), (key) => {
        return this.data[key] === null
      })
    },
    labels: function() {
      return Object.keys(this.data[this.section][0])
    },
    selectedSection: function() {
      return this.data[this.section]
    },
  },
  methods: {
    toHex: function(value) {
      return dec2HexString(value)
    },
    processData: function(value) {
      if (this._.isNumber(value)) return this.toHex(value)
      if (this._.isArray(value)) return this._.join(value, ", ")
      if (this._.isString(value)) return value
      if (this._.isNull(value)) return "none"
      return ""
    },
  },
}
</script>

<style lang="scss" scoped>
.container {
  overflow: hidden;
  margin-top: 2rem;

  .selector .button:hover,
  .selector .button.active {
    color: white;
    background-color: #00d1b2;
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
}
</style>
