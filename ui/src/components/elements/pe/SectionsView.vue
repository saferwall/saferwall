<template>
  <div class="container">
    <div class="sections_header">
      <div
        class="sections_header_field"
        v-for="(label, index) in labels"
        :key="index"
        :class="{ name: label === 'Name' }"
      >
        {{ _.startCase(label) }}
      </div>
    </div>
    <div class="section" v-for="(section, sec_index) in data" :key="sec_index">
      <div class="section_content">
        <div
          class="section_field"
          v-for="([key, value], third_index) in Object.entries(section)"
          :key="third_index"
          :class="{ name: key === 'Name' }"
        >
          <span class="parent">
            <span>
              {{ key == "Name" ? getName(value) : getHex(value) }}
            </span>
            <copy :content="key == 'Name' ? getName(value) : getHex(value)" />
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
  computed: {
    labels: function() {
      var keys = Object.keys(this.data[0])
      if (keys.length < 1) return []
      return keys
    },
  },
  methods: {
    getName: function(array) {
      var buffer = ""
      for (var dec of array) {
        buffer += String.fromCharCode(dec)
      }
      return buffer
    },
    getHex: function(dec) {
      return dec2HexString(dec)
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
    width: 10rem;
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
      width: 10rem;
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
