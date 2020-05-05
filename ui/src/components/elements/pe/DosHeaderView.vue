<template>
  <div class="container">
    <table class="table is-striped ">
      <thead>
        <tr>
          <th>Member</th>
          <th>Value</th>
          <th>Comment</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, index) in Object.entries(data)" :key="index">
          <th class="member">{{ _.startCase(item[0]) }}</th>
          <td class="value">
            <span class="parent">
              <span>
                {{ toHex(item[1]) }}
              </span>
              <copy :content="toHex(item[1])" />
            </span>
          </td>
          <td class="comment">
            {{ getComment(item[0], item[1]) }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { dec2HexString, reverse, hex2a, dec2Hex } from "@/helpers/pe"
import Copy from "@/components/elements/Copy"

export default {
  props: ["data"],
  components: {
    copy: Copy,
  },
  data() {
    return {}
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
        case "AddressOfNewEXEHeader":
          return reverse(hex2a(dec2Hex(value)))
        default:
          return ""
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.table {
  width: 100%;
  td {
    vertical-align: middle;
    &.value {
      .parent {
        position: relative;
        .copy {
          opacity: 0;
          transition: opacity 0.2s;
        }
        &:hover {
          .copy {
            opacity: 1;
          }
        }
      }
    }
  }
}
</style>
