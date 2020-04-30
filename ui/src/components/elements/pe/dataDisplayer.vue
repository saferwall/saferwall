<template>
  <div class="container" v-if="sectionName === 'DosHeader'">
    <table class="table is-striped is-narrow">
      <thead>
        <tr>
          <th>Member</th>
          <th>Value</th>
          <th>Comment</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, index) in Object.entries(data)" :key="index">
          <th class="member">{{ item[0] }}</th>
          <td class="value">
            {{ toHex(item[1]) }}
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

export default {
  props: ["data", "sectionName"],
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
          console.log(value)
          return reverse(hex2a(dec2Hex(value)))
        default:
          return ""
      }
    },
  },
}
</script>

<style lang="scss" scoped>
td {
  vertical-align: middle;
}
.member {
  width: max-content;
}
</style>
