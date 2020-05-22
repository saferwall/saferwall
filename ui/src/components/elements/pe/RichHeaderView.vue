<template>
  <div class="container">
    <div class="sections_header">
      <div
        class="sections_header_field"
        v-for="(label, index) in labels"
        :key="index"
      >
        {{ _.startCase(label) }}
      </div>
    </div>
    <div class="section" v-for="(section, sec_index) in data" :key="sec_index">
      <div class="section_content">
        <div class="section_field">
          <span class="parent">
            <span>
              {{ toHex(section["MinorCV"]) }}
            </span>
            <copy :content="toHex(section['MinorCV'])" />
          </span>
        </div>
        <div class="section_field">
          <span class="parent">
            <span>
              {{ toHex(section['ProdID']) }}
            </span>
            <copy :content="toHex(section['ProdID'])" />
          </span>
        </div>
        <div class="section_field">
          <span class="parent">
            <span>
              {{ section["Count"] }}
            </span>
            <copy :content="section['Count']" />
          </span>
        </div>
        <div class="section_field">
          <span class="parent">
            <span>
              {{ getMsInternalName(section['ProdID']) }}
            </span>
            <copy :content="getMsInternalName(section['ProdID'])" />
          </span>
        </div>
        <div class="section_field">
          <span class="parent">
            <span>
              {{ getVsVersion(section['ProdID']) }}
            </span>
            <copy :content="getVsVersion(section['ProdID'])" />
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { prodIdToStr, prodIdToVsVersion, dec2HexString } from "@/helpers/pe"
import Copy from "@/components/elements/Copy"

export default {
  props: ["data"],
  components: {
    copy: Copy,
  },
  data() {
    return {
      labels: [
        "MinorCV",
        "ProdId",
        "Count",
        "MS Internal Name",
        "Visual Studio Version",
      ],
    }
  },
  methods: {
    toHex: function(dec) {
      return dec2HexString(dec)
    },
    getMsInternalName: function(prodId) {
      return prodIdToStr(prodId)
    },
    getVsVersion: function(prodId) {
      return prodIdToVsVersion(prodId)
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
    width: 15rem;
    font-weight: 600;
  }
}
.section {
  padding: 0.2rem;
  .section_content {
    display: inline-flex;
    .section_field {
      text-align: left;
      width: 15rem;
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
