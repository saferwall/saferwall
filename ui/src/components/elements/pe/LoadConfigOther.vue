<template>
  <div class="container">
    <div class="sections_header">
      <div class="sections_header_field" v-for="label in labels" :key="label">
        {{ _.startCase(label) }}
      </div>
    </div>
    <div
      class="section"
      v-for="(content, index) in selectedSection"
      :key="index + 'b'"
    >
      <div class="section_content">
        <div
          class="section_field"
          v-for="(field, index) in Object.values(content)"
          :key="'function_' + index"
        >
          <span class="parent">
            <span>
              {{ processData(field) }}
              <copy :content="processData(field)" />
            </span>
            <span v-if="section === 'GFIDS' && index === 1">
              {{ getFlag(field) }}
              <copy :content="getFlag(field)" v-if="field" />
            </span>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Copy from "@/components/elements/Copy"
import { dec2HexString, GFIDS2String } from "@/helpers/pe"

export default {
  props: ["data", "section"],
  components: {
    copy: Copy,
  },
  computed: {
    labels: function() {
      if (this.section === "CHPE") return ["Structure Field", "Value"]
      if (this.section === "VolatileMetadata")
        return ["Structure Field", "Value", "Description"]
      if (this.section === "CFGLongJump") return ["RVA"]
      if (this.section === "SEH") return ["Handler"]
      if (this.section === "Access RVA") return ["Access RVA"]
      if (this.section === "Volatile Access Ranger") return ["RVA", "Size"]
      if (this.section === "Code Ranger") return ["Begin", "End", "Machine"]
      return Object.keys(this.data[this.section][0])
    },
    selectedSection: function() {
      if (this.section === "CHPE")
        return this._.toPairs(this.data[this.section].CHPEMetadata)
      if (this.section === "Volatile Access Ranger")
        return this.data.VolatileMetadata.InfoRangeTable
      if (this.section === "Code Ranger")
        return this.data.CHPE.CodeRanges
      if (this.section === "Access RVA")
        return this._.map(this.data.VolatileMetadata.AccessRVATable, (a) => {
          return { a }
        })
      if (this.section === "VolatileMetadata") {
        var data = this._.toPairs(this.data[this.section].Struct)
        return data.map((item) => {
          if (
            item[0] === "Size" ||
            item[0] === "VolatileAccessTableSize" ||
            item[0] === "VolatileInfoRangeTableSize"
          ) {
            item.push(this.getSize(item[1]))
          }
          return item
        })
      }
      if (this.section === "CFGLongJump" || this.section === "SEH")
        return this._.map(this.data[this.section], (a) => {
          return { a }
        })
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
    getFlag: function(value) {
      if (value === 0) return ""
      return GFIDS2String(value)
    },
    getSize: function(value) {
      if (value >= 1000000) return (value / 1000000).toFixed(2) + " MB"
      if (value >= 1000) return (value / 1000).toFixed(2) + " KB"
      else return value + " B"
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
      width: 25rem;
      font-weight: 600;
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
