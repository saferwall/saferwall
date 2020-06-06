<template>
  <div class="container" v-if="this.data">
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
  data() {
    return {
      labels: ["Ordinal", "FunctionRVA", "Name"],
    }
  },
  computed: {
    filteredData: function() {
      return this._.map(this.data, (data) => {
        return {
          Ordinal: data.Ordinal,
          FunctionRVA: data.FunctionRVA,
          Name: data.Name,
        }
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
  }
}
.section {
  padding: 0.2rem;
  .section_content {
    display: inline-flex;
    .section_field {
      text-align: left;
      width: 17rem;
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
