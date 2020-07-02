<template>
  <div v-if="isPE">
    <div class="tile is-vertical box">
      <div class="title">Portable Executable</div>
      <div class="columns">
        <div class="column is-one-quarter items">
          <div
            class="bt"
            :class="{
              active: selectedSection === section,
            }"
            v-for="(section, index) in sections"
            :key="index"
          >
            <button class="button is-light" @click="selectedSection = section" :class="{subsection: loadConfigSubSections.includes(section)}">
              {{ _.startCase(section) }}
            </button>
            <hr />
          </div>
        </div>
        <div class="column values">
          <DosHeaderView
            v-if="selectedSection === 'DosHeader'"
            :data="pe[selectedSection]"
            :signature="pe['NtHeader'].Signature"
          />
          <RichHeaderView
            v-if="selectedSection === 'RichHeader'"
            :data="pe[selectedSection].CompIDs"
          />
          <NtHeaderView
            v-if="selectedSection === 'NtHeader'"
            :data="pe[selectedSection]"
          />
          <SectionsView
            v-if="selectedSection === 'Sections'"
            :data="pe[selectedSection]"
          />
          <Imports
            v-if="selectedSection === 'Imports'"
            :data="pe[selectedSection]"
          />
          <Export
            v-if="selectedSection === 'Export'"
            :data="pe[selectedSection]"
          />
          <LoadConfig
            v-if="selectedSection === 'LoadConfig'"
            :data="pe[selectedSection]"
          />
          <LoadConfigOther
            v-if="loadConfigSubSections.includes(selectedSection)"
            :data="_.omit(pe['LoadConfig'], 'LoadCfgStruct')"
            :section="selectedSection"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import DosHeaderView from "../elements/pe/DosHeaderView"
import RichHeaderView from "../elements/pe/RichHeaderView"
import NtHeaderView from "../elements/pe/NtHeaderView"
import SectionsView from "../elements/pe/SectionsView"
import Imports from "../elements/pe/Imports"
import Export from "../elements/pe/Export"
import LoadConfig from "../elements/pe/LoadConfig"
import LoadConfigOther from "../elements/pe/LoadConfigOther"
import { mapGetters } from "vuex"

export default {
  components: {
    DosHeaderView,
    RichHeaderView,
    NtHeaderView,
    SectionsView,
    Imports,
    Export,
    LoadConfig,
    LoadConfigOther,
  },
  data() {
    return {
      selectedSection: "DosHeader",
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
      var keys = Object.keys(this.pe)
      var index = keys.indexOf("LoadConfig")
      return this._.concat(
        this._.slice(keys, 0, index + 1),
        this.loadConfigSubSections,
        this._.slice(keys, index + 1, keys.length),
      )
    },
    loadConfigSubSections: function() {
      var sections = Object.keys(this.pe.LoadConfig)
      return sections.filter((section) => {
        return (
          section !== "LoadCfgStruct" && this.pe.LoadConfig[section] !== null
        )
      })
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
          justify-content: left;
          width: 10em;
        }
        .subsection{
          padding-left: 2rem;
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
