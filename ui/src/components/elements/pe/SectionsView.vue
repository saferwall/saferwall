<template>
  <div class="container">
    <div class="sections_header">
      <div
        class="sections_header_field"
        :class="{
          Name: label === 'Name',
          VirtualSize: label === 'VirtualSize',
          VirtualAddress: label === 'VirtualAddress',
          SizeOfRawData: label === 'SizeOfRawData',
          NumberOfLineNumbers: label === 'NumberOfLineNumbers',
          PointerToRawData: label === 'PointerToRawData',
        }"
        v-for="(label, index) in labels"
        :key="index"
      >
        {{ label }}
      </div>
    </div>
    <div class="section" v-for="(section, sec_index) in data" :key="sec_index">
      <div class="section_content">
        <div
          class="section_field"
          :class="{
            Name: key === 'Name',
            VirtualSize: key === 'VirtualSize',
            VirtualAddress: key === 'VirtualAddress',
            SizeOfRawData: key === 'SizeOfRawData',
            NumberOfLineNumbers: key === 'NumberOfLineNumbers',
            PointerToRawData: key === 'PointerToRawData',
          }"
          v-for="([key, value], third_index) in Object.entries(section)"
          :key="third_index"
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
  padding: 0.5rem;
  .sections_header_field {
    text-align: left;
    width: 10rem;
    margin-right: 1rem;
    font-weight: 600;
  }
}
.section {
  padding: 0.5rem;
  .section_content {
    display: inline-flex;
    .section_field {
      text-align: left;
      width: 10rem;
      margin-right: 1rem;
      &:hover {
        .copy {
          opacity: 1;
        }
      }
    }
  }
}
.Name,
.VirtualSize {
  width: 6rem !important;
}
.VirtualAddress {
  width: 7rem !important;
}
.SizeOfRawData {
  width: 8rem !important;
}
.NumberOfLineNumbers {
  margin-right: 1.5rem !important;
}
.PointerToRawData {
  width: 9rem !important;
}
.parent {
  position: relative;
  .copy {
    opacity: 0;
    transition: opacity 0.2s;
  }
}
</style>
