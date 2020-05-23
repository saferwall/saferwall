<template>
  <div class="container">
    <div class="sections_header">
      <div
        class="sections_header_field member"
      >
        {{ _.startCase(labels[0]) }}
      </div>
      <div
        class="sections_header_field value"
      >
        {{ _.startCase(labels[1]) }}
      </div>
      <div
        class="sections_header_field comment"
      >
        {{ _.startCase(labels[2]) }}
      </div>
    </div>
    <div
      class="section"
      v-for="(section, sec_index) in Object.entries(data)"
      :key="sec_index"
    >
      <div class="section_content">
        <div class="section_field member">
          <span class="parent">
            <span>
              {{ _.startCase(section[0]) }}
            </span>
          </span>
        </div>
        <div class="section_field value">
          <span class="parent">
            <span>
              {{ toHex(section[1]) }}
            </span>
            <copy :content="toHex(section[1])" />
          </span>
        </div>
        <div class="section_field comment">
          <span class="parent">
            <span>
              {{ getComment(section[0], section[1]) }}
            </span>
            <copy :content="getComment(section[0], section[1])" />
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { dec2HexString, reverse, hex2a, dec2Hex } from "@/helpers/pe"
import Copy from "@/components/elements/Copy"

export default {
  props: ["data", "signature"],
  components: {
    copy: Copy,
  },
  data() {
    return {
      labels: ["Member", "Value", "Comment"],
    }
  },
  methods: {
    toHex: function(value) {
      if (Array.isArray(value)) {
        var tmpArray = []
        for (var index in value) {
          tmpArray.push(dec2HexString(value[index]))
        }
        return tmpArray
      } else return dec2HexString(value)
    },
    getComment: function(name, value) {
      switch (name) {
        case "Magic":
          return reverse(hex2a(dec2Hex(value)))
        case "AddressOfNewEXEHeader":
          return reverse(hex2a(dec2Hex(this.signature)))
        default:
          return ""
      }
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
    margin-right: 1rem;
    font-weight: 600;
    &.member {
      width: 25rem;
    }
    &.value {
      width: 35rem;
    }
  }
}
.section {
  padding: 0.2rem;
  .section_content {
    display: inline-flex;
    .section_field {
      text-align: left;
      margin-right: 1rem;
      &.member {
        width: 25rem;
        font-weight: 600;
      }
      &.value {
        width: 35rem;
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
