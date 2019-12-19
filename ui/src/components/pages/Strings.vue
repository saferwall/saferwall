<!-- refactor (filtering, sorting…) -->

<template>
  <div>
    <loader v-if="showLoader"></loader>
    <div class="columns">
      <div class="column is-6 is-offset-6" style="text-align:right">
        <label for="" class="label">Limit:</label>
        <div class="select is-medium">
          <select id="select-limit" v-model="limit" @change="limitChanged()">
            <option value="10">10</option>
            <option value="100">100</option>
            <option value="1000">1000</option>
          </select>
        </div>
      </div>
    </div>
    <div class="tile is-ancestor" v-if="!showLoader">
      <div class="tile is-parent">
        <div class="tile is-child box">
          <table class="table is-bordered is-striped">
            <thead
              style="display:block;position:sticky;top:50px;background: #fff"
            >
              <tr>
                <th>
                  <div class="head-title">Encoding</div>
                  <div class="sort-icons">
                    <i
                      class="icon ion-arrow-up-b up"
                      @click="sortTable('asc', 'encoding')"
                    ></i>
                    <i
                      class="icon ion-arrow-down-b down"
                      @click="sortTable('desc', 'encoding')"
                    ></i>
                  </div>
                </th>
                <th>
                  <span class="head-title">Value</span>
                  <div class="sort-icons">
                    <i
                      class="icon ion-arrow-up-b up"
                      @click="sortTable('asc', 'value')"
                    ></i>
                    <i
                      class="icon ion-arrow-down-b down"
                      @click="sortTable('desc', 'value')"
                    ></i>
                  </div>
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(string, index) in filteredStringsToShow" :key="index">
                <td>{{ string.encoding }}</td>
                <td>{{ string.value }}</td>
              </tr>
              <tr>
                <td>
                  <div class="control">
                    <input
                      class="input"
                      @keyup="search('encodings')"
                      v-model="encodingSearched"
                      type="text"
                      placeholder="Search Encodings..."
                    />
                  </div>
                </td>
                <td>
                  <div class="control">
                    <input
                      class="input"
                      @keyup="search('values')"
                      v-model="valueSearched"
                      type="text"
                      placeholder="Search Values..."
                    />
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import Loader from "@/components/elements/Loader"

export default {
  components: {
    loader: Loader,
  },
  data() {
    return {
      showLoader: true,
      strings: [],
      filteredStrings: [],
      start: 0,
      limit: 10,
      encodingSorted: "",
      valueSorted: "",
      valueSearched: "",
      encodingSearched: "",
    }
  },
  computed: {
    filteredStringsToShow: function() {
      return this.filteredStrings.filter((string) => string.show)
    },
  },
  mounted() {
    if (this.$store.getters.getHashContext) this.showData(this.$store.getters.getHashContext)
  },

  methods: {
    sortTable(key, column) {
      if (column === "encoding") {
        if (key === "asc") {
          this.filteredStrings.sort((a, b) =>
            a.encoding > b.encoding ? 1 : b.encoding > a.encoding ? -1 : 0,
          )
        }
        if (key === "desc") {
          this.filteredStrings.sort((a, b) =>
            a.encoding > b.encoding ? -1 : b.encoding > a.encoding ? 1 : 0,
          )
        }

        this.filteredStrings = this.filteredStrings
          .filter((s) => s.show)
          .slice(this.start, this.limit)
      } else if (column === "value") {
        if (key === "asc") {
          this.filteredStrings.sort((a, b) =>
            a.value > b.value ? 1 : b.value > a.value ? -1 : 0,
          )
        }
        if (key === "desc") {
          this.filteredStrings.sort((a, b) =>
            a.value > b.value ? -1 : b.value > a.value ? 1 : 0,
          )
        }

        this.filteredStrings = this.filteredStrings
          .filter((s) => s.show)
          .slice(this.start, this.limit)
      }
    },
    limitChanged() {
      this.filteredStrings = this.strings.slice(this.start, this.limit)
    },
    search(key) {
      if (this.encodingSearched.length > 0 || this.valueSearched.length > 0) {
        if (key === "encodings") {
          this.strings.forEach(
            (s) =>
              (s.show = s.encoding
                .trim()
                .toLowerCase()
                .startsWith(this.encodingSearched.toLowerCase())),
          )
          this.filteredStrings = this.strings
            .filter((s) => s.show)
            .slice(this.start, this.limit)
        }
        if (key === "values") {
          this.strings.forEach(
            (s) =>
              (s.show = s.value
                .trim()
                .toLowerCase()
                .startsWith(this.valueSearched.toLowerCase())),
          )
          this.filteredStrings = this.strings
            .filter((s) => s.show)
            .slice(this.start, this.limit)
          console.log(this.filteredStrings)
        }
      } else {
        this.strings.forEach((s) => (s.show = true))
        this.filteredStrings = this.strings.slice(this.start, this.limit)
      }
    },
    showData(hash) {
      this.$http
        .get(this.$api_endpoints.GET_FILES + hash)
        .then((data) => {
          this.showLoader = false
          this.strings = data.data.strings
          this.strings.forEach((s) => (s.show = true))
          this.filteredStrings = this.strings.slice(this.start, this.limit)
        })
        .catch((err) => console.error(err))
    },
  },
}
</script>
<style lang="scss" scoped>
.head-title {
  float: left;
  line-height: 30px;
}

tbody,
thead tr {
  display: table;
  width: 100%;
  table-layout: fixed;
}

.sort-icons {
  float: left;
  position: relative;
  height: 30px;
  margin-left: 10px;
  margin-right: 10px;

  .icon {
    position: absolute;
    cursor: pointer;
    opacity: 0.25;
    transition: all 0.2s;
    height: auto;
    width: auto;
    line-height: normal;

    &:hover,
    &.active,
    &:active {
      opacity: 1;
    }

    &.up {
      top: calc(50% - 2px);
      transform: translateY(-50%);
    }

    &.down {
      top: 50%;
    }
  }
}
</style>
