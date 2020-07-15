<template>
  <div class="container">
    <div class="sections_header">
      <div
        class="sections_header_field"
        v-for="(label, index) in labels"
        :key="index"
        :class="{ name: label === 'Name' }"
      >
        {{ label }}
      </div>
    </div>
    <div class="section" v-for="(section, sec_index) in data" :key="sec_index">
      <div class="section_content">
        <div
          class="section_field"
          v-for="([key, value], third_index) in Object.entries(section)"
          :key="third_index"
        >
          <span class="parent">
            <span>
              {{ processData(value, third_index) }}
            </span>
            <copy :content="processData(value, third_index)" />
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { dec2HexString } from "../../../helpers/pe"
import Copy from "@/components/elements/Copy"

export default {
  props: ["data"],
  components: {
    copy: Copy,
  },
  data() {
    return {
      labels: ['#', 'Entry RVA', 'Value', 'Description']
    }
  },
  methods: {
    toHex: function(dec) {
      return dec2HexString(dec)
    },
    processData: function(value, index) {
      if(index === 0) return value
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
}
.sections_header {
  display: inline-flex;
  padding: 0.2rem;
  .sections_header_field {
    text-align: left;
    width: 11rem;
    font-weight: 600;
  }
}
.section {
  padding: 0.2rem;
  .section_content {
    display: inline-flex;
    .section_field {
      text-align: left;
      width: 11rem;
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
