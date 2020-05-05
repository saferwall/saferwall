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
            }"
            v-for="(section, index) in sections"
            :key="index"
          >
            <button class="button is-light" @click="selectedSection = index">
              {{ _.startCase(section) }}
            </button>
            <hr />
          </div>
        </div>
        <div class="column values">
          <DosHeaderView
            v-if="sections[selectedSection] === 'DosHeader'"
            :data="pe[sections[selectedSection]]"
          />
          <RichHeaderView 
          v-if="sections[selectedSection] === 'RichHeader'"
            :data="pe[sections[selectedSection]].CompIDs"/>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import DosHeaderView from "../elements/pe/DosHeaderView"
import RichHeaderView from "../elements/pe/RichHeaderView"
import { mapGetters } from "vuex"

export default {
  components: {
    DosHeaderView,
    RichHeaderView,
  },
  data() {
    return {
      selectedSection: 1,
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
      .bt {
        display: inline-flex;
        height: 35px;
        margin-bottom: 0;
        margin-top: 0;

        .button {
          background-color: transparent;
          font-size: medium;
          font-weight: 400;
          width: 300px;
          justify-content: right;
          width: 10em;
        }
        hr {
          background-color: grey;
          border: none;
          display: block;
          height: 35px;
          width: 1px;
          margin: 0;
        }
      }

      .active {
        button {
          color: #00d1b2;
          font-weight: 600;
        }
        hr {
          background-color: #00d1b2;
          width: 1.5px;
        }
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
    .values {
      margin-left: 2rem;
    }
  }
}
</style>
