<template>
  <div class="tile is-child box">
    <h4 class="title">Submissions</h4>
    <table class="table is-striped data-data">
      <thead>
        <tr>
          <th>Data</th>
          <th>Name</th>
          <th>Source</th>
          <th>Country</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(i, index) in summaryData.submissions" :key="index">
          <td>
            <div class="data-value">
              <span class="value-text">{{ formateDate(i.date) }}</span>
              <copy :content="formateDate(i.date)"></copy>
            </div>
          </td>
          <td>
            <div class="data-value">
              <span class="value-text">{{ i.filename }}</span>
              <copy :content="i.filename"></copy>
            </div>
          </td>
          <td>
            <div class="data-value">
              <span
                ><i
                  :class="'icon mdi mdi-24px mdi-' + source_icons[i.source]"
                ></i
              ></span>
              <span class="value-text">&nbsp;{{ i.source }}</span>
              <copy :content="i.source"></copy>
            </div>
          </td>
          <td>
            <div class="data-value">
              <span
                :class="'flag-icon flag-icon-' + i.country.toLowerCase()"
              ></span>
              <span class="value-text">&nbsp;{{ i.country }}</span>
              <copy :content="i.country"></copy>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import Copy from "@/components/elements/Copy"

export default {
  props: ["summaryData"],
  components: {
    copy: Copy,
  },
  data() {
    return {
      source_icons: {
        api: "server",
        web: "earth",
      },
    }
  },
  methods: {
    formateDate: function(rawDate) {
      var elements = rawDate.split("T")
      return elements[0] + " " + elements[1].split(".")[0]
    },
  },
}
</script>

<style lang="scss" scoped>
@import "flag-icon-css/css/flag-icon.min.css";
@import url("https://cdn.jsdelivr.net/npm/@mdi/font@4.x/css/materialdesignicons.min.css");

.data-data.table {
  background: white;
  .value-text {
    vertical-align: top;
  }
}
.tile {
  overflow-y: auto;
  height: 45em;
}
</style>
