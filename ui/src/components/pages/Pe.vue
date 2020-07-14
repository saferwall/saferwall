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
            <button
              class="button is-light"
              @click="selectedSection = section"
              :class="{ subsection: loadConfigSubSections.includes(section), extrasubsection: extraSubSections.includes(section) }"
            >
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
          <TLS
            v-if="selectedSection === 'TLS'"
            :data="pe[selectedSection]"
          />
          <Debugs
            v-if="selectedSection === 'Debugs'"
            :data="pe[selectedSection]"
          />
          <IAT
            v-if="selectedSection === 'IAT'"
            :data="pe[selectedSection]"
          />
          <LoadConfig
            v-if="selectedSection === 'LoadConfig'"
            :data="pe[selectedSection]"
          />
          <Enclave
            v-if="selectedSection === 'Enclave'"
            :data="pe['LoadConfig']['Enclave']"
          />
          <LoadConfigOther
            v-if="
              loadConfigSubSections.includes(selectedSection) &&
                selectedSection !== 'Enclave'
            "
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
import TLS from "../elements/pe/TLS"
import Debugs from "../elements/pe/Debugs"
import IAT from "../elements/pe/IAT"
import LoadConfig from "../elements/pe/LoadConfig"
import LoadConfigOther from "../elements/pe/LoadConfigOther"
import Enclave from "../elements/pe/Enclave"
import { mapGetters } from "vuex"

export default {
  components: {
    DosHeaderView,
    RichHeaderView,
    NtHeaderView,
    SectionsView,
    Imports,
    Export,
    TLS,
    Debugs,
    IAT,
    LoadConfig,
    LoadConfigOther,
    Enclave,
  },
  data() {
    return {
      selectedSection: "DosHeader",
      extraSubSections : ['Access RVA','Volatile Access Ranger', 'Code Ranger']
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
      keys = keys.filter((section) => this.pe[section] !== null)
      return this.addAtPosition(keys, this.loadConfigSubSections, "LoadConfig")
    },
    loadConfigSubSections: function() {
      var sections = Object.keys(this.pe.LoadConfig)
      sections = sections.filter((section) => {
        return (
          section !== "LoadCfgStruct" && this.pe.LoadConfig[section] !== null
        )
      })
      sections = this.addAtPosition(sections, ['Access RVA','Volatile Access Ranger'], 'VolatileMetadata')
      sections = this.addAtPosition(sections, ['Code Ranger'], 'CHPE')

      return sections
    },
  },
  methods: {
    addAtPosition: function(array, element, position) {
      var index = this._.indexOf(array, position)
      return this._.concat(
        this._.slice(array, 0, index + 1),
        element,
        this._.slice(array, index + 1, array.length),
      )
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
    padding: 0.5rem;
    .items {
      display: grid;
      width: 20rem;
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
        .subsection {
          padding-left: 2rem;
        }
        .extrasubsection {
          padding-left: 3rem;
        }
        hr {
          background-color: grey;
          border: none;
          display: block;
          height: 35px;
          width: 1px;
          margin: auto !important;
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
  }
}
</style>
