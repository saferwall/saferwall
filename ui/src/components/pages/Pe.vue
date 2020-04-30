<template>
  <div v-if="isPE">
    <div class="tile is-vertical box">
      <div class="title">Portable Executable</div>
      <div class="columns">
        <div class="column is-one-quarter items">
          <div
            class="bt"
            :class="{
              active: selectedSection === index,
              grey: index % 2 === 0,
            }"
            v-for="(section, index) in sections"
            :key="index"
          >
            <button class="button is-light" @click="selectedSection = index">
              {{ section }}
            </button>
            <hr />
          </div>
        </div>
        <div class="column values">
          <dataDisplayer :data="pe[sections[selectedSection]]" :sectionName="sections[selectedSection]"/>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import dataDisplayer from "../elements/pe/dataDisplayer"
import { mapGetters } from "vuex"

export default {
  components:{
    dataDisplayer,
  },
  data() {
    return {
      selectedSection: 0,
    }
  },
  computed: {
    ...mapGetters({ fileData: "getFileData", isPE: "isPE" }),
    pe: function() {
      if (
        !this.fileData ||
        (Object.entries(this.fileData).length === 0 &&
          this.fileData.constructor === Object)
      )
        return {}
      return this._.omit(this.fileData.pe, [
        "Certificates",
        "BoundImports",
        "GlobalPtr",
        "CLRHeader",
        "Header",
        "Is64",
        "Is32",
        "Anomalies",
      ])
    },
    sections: function() {
      return Object.keys(this.pe)
    },
  },
}
</script>

<style lang="scss" scoped>
.tile {
  padding-bottom: 1em;
  width: 60%;
  margin: auto;
  .title {
    margin-bottom: 20px;
  }
  .columns {
    padding: 1rem;
    .items {
      display: grid;
      width: max-content;
      height: max-content;
      text-align: left;

      .bt {
        display: inline-flex;
        height: 35px;
        margin-bottom: 0;
        margin-top: 0;

        .button {
          background-color: transparent;
          width: 200px;
          justify-content: left;
          width: 10em;
        }
        hr {
          background-color: grey;
          border: none;
          display: block;
          height: 35px;
          width: 2px;
          margin: 0;
        }
      }

      .active {
        background-color: rgba(0, 215, 211, 0.1) !important;
        button {
          color: #00d1b2;
          font-weight: 600;
        }
        hr {
          background-color: #00d1b2;
        }
      }

      .grey {
        background-color: #f7f7f7;
      }

      .bt:hover {
        button {
          color: #00d1b2;
          font-weight: 600;
        }
        hr {
          background-color: #00d1b2;
        }
      }
    }
    .values{
      margin-left: 2rem;
    }
  }
}
</style>
