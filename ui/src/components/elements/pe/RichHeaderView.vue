<template>
  <div class="container">
    <table class="table is-striped ">
      <thead>
        <tr>
          <th>MinorCV</th>
          <th>ProdId</th>
          <th>Count</th>
          <th>MS Internal Name</th>
          <th>Visual Studio Version</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, index) in data" :key="index">
          <td class="MinorCV">
            <span class="parent">
              <span>
                {{ toHex(item["MinorCV"]) }}
              </span>
              <copy :content="toHex(item['MinorCV'])" />
            </span>
          </td>
          <td class="ProdId">
            <span class="parent">
              <span>
                {{ toHex(item["ProdID"]) }}
              </span>
              <copy :content="toHex(item['ProdID'])" />
            </span>
          </td>
          <td class="Count">
            <span class="parent">
              <span>
                {{ item["Count"] }}
              </span>
              <copy :content="item['Count']" />
            </span>
          </td>
          <td class="MSInternalName">
            <span class="parent">
              <span>
                {{ getMsInternalName(item["ProdID"]) }}
              </span>
              <copy :content="getMsInternalName(item['ProdID'])" />
            </span>
          </td>
          <td class="VSVersion">
            <span class="parent">
              <span>
                {{ getVsVersion(item["ProdID"]) }}
              </span>
              <copy :content="getVsVersion(item['ProdID'])" />
            </span>
          </td>
        </tr>
      </tbody>
    </table>
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
.table {
  width: 100%;
  td {
    vertical-align: middle;
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
</style>
