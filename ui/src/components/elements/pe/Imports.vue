<template>
  <div class="container">
    <div class="sections_header">
      <div
        class="sections_header_field"
        v-for="label in labels"
        :key="label"
      >
        {{ _.startCase(label) }}
      </div>
    </div>
    <div
      class="section"
      v-for="(section, index) in data"
      :key="index"
      @click="select(index)"
    >
      <div class="section_content">
        <div class="section_field">
          <span class="parent">
            <span>
              {{ section["Name"] }}
            </span>
            <copy :content="section['Name']" />
          </span>
        </div>
        <div class="section_field">
          <span class="parent">
            <span>
              {{ toHex(section["Descriptor"]["OriginalFirstThunk"]) }}
            </span>
            <copy
              :content="toHex(section['Descriptor']['OriginalFirstThunk'])"
            />
          </span>
        </div>
        <div class="section_field">
          <span class="parent">
            <span>
              {{ toHex(section["Descriptor"]["TimeDateStamp"]) }}
            </span>
            <copy :content="toHex(section['Descriptor']['TimeDateStamp'])" />
          </span>
        </div>
        <div class="section_field">
          <span class="parent">
            <span>
              {{ toHex(section["Descriptor"]["ForwarderChain"]) }}
            </span>
            <copy :content="toHex(section['Descriptor']['ForwarderChain'])" />
          </span>
        </div>
        <div class="section_field">
          <span class="parent">
            <span>
              {{ toHex(section["Descriptor"]["Name"]) }}
            </span>
            <copy :content="toHex(section['Descriptor']['Name'])" />
          </span>
        </div>
        <div class="section_field">
          <span class="parent">
            <span>
              {{ toHex(section["Descriptor"]["FirstThunk"]) }}
            </span>
            <copy :content="toHex(section['Descriptor']['FirstThunk'])" />
          </span>
        </div>
      </div>
    </div>
    <ImportsFunctions :data="selectedFunction" :key="selectedIndex+'_functions'" />
  </div>
</template>

<script>
import { dec2HexString } from "@/helpers/pe"
import Copy from "@/components/elements/Copy"
import ImportsFunctions from "./Imports_functions"

export default {
  props: ["data"],
  components: {
    copy: Copy,
    ImportsFunctions,
  },
  data() {
    return {
      selectedIndex: -1,
    }
  },
  computed: {
    labels: function() {
      if (this.data)
        return this._.concat(["NameRVA"], Object.keys(this.data[0].Descriptor))
      return null
    },
    selectedFunction: function() {
      return this.selectedIndex > -1
        ? this.data[this.selectedIndex].Functions
        : null
    },
  },
  methods: {
    toHex: function(value) {
      return dec2HexString(value)
    },
    select: function(index) {
      this.selectedIndex = index
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
    width: 17rem;
    font-weight: 600;
    &.name {
      width: 6rem;
    }
  }
}
.section {
  padding: 0.2rem;
  &:hover {
    background-color: #f7f7f7;
    cursor: pointer;
  }
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
