<template>
  <div class="tile is-child box">
    <h4 class="title">Submissions</h4>
    <table class="table is-fullwidth is-striped data-data">
      <thead>
        <tr>
          <th>Data</th>
          <th>Name</th>
          <th>Source</th>
          <th>Country</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(i, index) in submissions" :key="index" >
          <td>
            <div class="data-value">
              <span class="value-text">{{ i.date }}</span>
              <copy :content="i.date"></copy>
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
              <span><i :class="'icon mdi mdi-24px mdi-'+source_icons[i.source]"></i></span>
              <span class="value-text">&nbsp;{{ i.source }}</span>
              <copy :content="i.source"></copy>
            </div>
          </td>
          <td>
            <div class="data-value">
              <span :class="'flag-icon flag-icon-'+i.country.toLowerCase()"></span>
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
      source_icons:{
        api : "server",
        web : "earth",
      }
    }
  },
  computed: {
    submissions: function() {
      return this.summaryData.submissions
    },
  },
}
</script>

<style lang="scss" scoped>
@import 'flag-icon-css/css/flag-icon.min.css';
@import url("https://cdn.jsdelivr.net/npm/@mdi/font@4.x/css/materialdesignicons.min.css");

.data-data.table{
  background: white;
  .value-text{
    vertical-align: top;
  }
}
</style>
