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
      v-for="(section, sec_index) in Object.entries(LoadCfgStruct)"
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
          <span v-if="section[0] === 'GuardFlags'">
            <span class="parent">
              <span
                class="tag is-link is-normal"
                :class="{ redTag: isAntivirusTag(tag[0]) }"
                id="tag"
                v-for="tag in tags"
                :key="tag"
              >
                {{ tag }}
              </span>
            </span>
          </span>
          <span class="parent" v-else>
            <span>
              {{ getComment(section[0], section[1]) }}
            </span>
            <copy
              v-if="getComment(section[0], section[1])"
              :content="getComment(section[0], section[1])"
            />
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { dec2HexString, GuardFlags2String } from "@/helpers/pe"
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
    LoadCfgStruct: function() {
      var data = this._.omit(this.data.LoadCfgStruct, "CodeIntegrity")
      data[
        "CodeIntegrity.Catalog"
      ] = this.data.LoadCfgStruct.CodeIntegrity.Catalog
      data[
        "CodeIntegrity.CatalogOffset"
      ] = this.data.LoadCfgStruct.CodeIntegrity.CatalogOffset
      data["CodeIntegrity.Flags"] = this.data.LoadCfgStruct.CodeIntegrity.Flags
      data[
        "CodeIntegrity.Reserved"
      ] = this.data.LoadCfgStruct.CodeIntegrity.Reserved

      return data
    },
    tags: function() {
      if (this.data)
        return GuardFlags2String(this.data.LoadCfgStruct.GuardFlags)
      return []
    },
  },
  methods: {
    toHex: function(value) {
      return dec2HexString(value)
    },
    getComment: function(field, value) {
      switch (field) {
        default:
          return ""
      }
    },
    isAntivirusTag: function(tag) {
      const antivirusList = [
        "eset",
        "fsecure",
        "avira",
        "bitdefender",
        "kaspersky",
        "symantec",
        "sophos",
        "windefender",
        "clamav",
        "comodo",
        "avast",
        "mcafee",
      ]
      return antivirusList.includes(tag)
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
      width: 10rem;
    }
    &.comment {
      width: 50rem;
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
        font-weight: 500;
      }
      &.value {
        width: 10rem;
      }
      &.comment {
        width: 50rem;
        #tag {
          margin-right: 0.2em;
          color: white;
          font-weight: 600;
        }
        .redTag {
          background-color: #f14668;
        }
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
