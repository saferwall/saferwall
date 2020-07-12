<template>
  <div class="container">
    <div class="sections_header">
      <div class="sections_header_field member">
        {{ _.startCase(labels[0]) }}
      </div>
      <div class="sections_header_field value">
        {{ _.startCase(labels[1]) }}
      </div>
      <div class="sections_header_field comment">
        {{ _.startCase(labels[2]) }}
      </div>
    </div>
    <div
      class="section"
      v-for="(section, sec_index) in Object.entries(data.Config)"
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
              {{ getValue(section) }}
            </span>
            <copy :content="getValue(section)" />
          </span>
        </div>
        <div class="section_field comment">
          <span class="parent">
            <span>
              {{ getDescription(section) }}
            </span>
            <copy
              v-if="getDescription(section)"
              :content="getDescription(section)"
            />
          </span>
        </div>
      </div>
    </div>
    <div class="SecondSection">
      <div class="sections_header second">
        <div
          class="sections_header_field second"
          v-for="(label, index) in secondSectionLabels"
          :key="index"
        >
          {{ _.startCase(label) }}
        </div>
      </div>
      <div
        class="section"
        v-for="(section, index) in data.Imports"
        :key="index"
      >
        <div class="section_content">
          <div
            class="section_field second"
            v-for="(row, index) in Object.values(section)"
            :key="index"
          >
            <span class="parent">
              <span>
                {{ _.truncate(toHex(row)) }}
                <copy :content="toHex(row)" />
              </span>
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { dec2HexString, dec2Hex } from "@/helpers/pe"
import Copy from "@/components/elements/Copy"

export default {
  props: ["data"],
  components: {
    copy: Copy,
  },
  data() {
    return {
      labels: ["Structure Field", "Value", "Description"],
    }
  },
  computed: {
    secondSectionLabels: function() {
      if (this.data) return Object.keys(this.data.Imports[0])
      return []
    },
  },
  methods: {
    toHex: function(value) {
      if (Array.isArray(value)) {
        var tmpArray = []
        for (var index in value) {
          tmpArray.push(dec2Hex(value[index]))
        }
        return this._.join(tmpArray, " ")
      }
      return dec2HexString(value)
    },
    getSize: function(value) {
      if (value >= 1000000) return (value / 1000000).toFixed(2) + " MB"
      if (value >= 1000) return (value / 1000).toFixed(2) + " KB"
      else return value + " B"
    },
    getValue: function(sec) {
      switch (sec[0]) {
        case "FamilyID":
        case "ImageID":
          return ""
        default:
          return this.toHex(sec[1])
      }
    },
    getDescription: function(sec) {
      switch (sec[0]) {
        case "Size":
        case "MinimumRequiredConfigSize":
          return this.getSize(sec[1])
        case "FamilyID":
        case "ImageID":
          return this.toHex(sec[1])
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
      width: 24rem;
    }
    &.value {
      width: 24rem;
    }
    &.second {
      width: 13rem;
    }
  }
}
.section {
  padding: 0.2rem;
  .section_content {
    display: inline-flex;
    .section_field {
      text-align: left;
      width: 25rem;
      &:hover {
        .copy {
          opacity: 1;
        }
      }
    }
    .section_field.second {
      width: 14rem;
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
.SecondSection {
  margin-top: 2rem !important;
}
</style>
