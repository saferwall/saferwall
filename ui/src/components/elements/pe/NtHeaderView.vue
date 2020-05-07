<template>
  <div class="container">
    <div
      class="section"
      v-for="(section, index) in Object.entries(sections)"
      :key="index"
    >
      <div class="section_title">{{ section[0] }}</div>
      <div
        class="section_content"
        v-for="(line, index) in section[1]"
        :key="index"
      >
        <div class="multiple" v-if="line.type === 'multiple'">
          <div class="first_field">
            <div class="section_field_name">{{ Object.keys(line)[1] }}</div>
            <div class="section_field_value">
              {{ toHex(Object.values(line)[1]) }}
            </div>
            <div class="section_field_description">
              {{ getDescription(Object.keys(line)[1], Object.values(line)[1]) }}
            </div>
          </div>
          <div class="second_field">
            <div class="section_field_name">{{ Object.keys(line)[2] }}</div>
            <div class="section_field_value">
              {{ toHex(Object.values(line)[2]) }}
            </div>
            <div class="section_field_description">
              {{ getDescription(Object.keys(line)[2], Object.values(line)[2]) }}
            </div>
          </div>
        </div>
        <div class="single" v-else>
          <div class="section_field_name">{{ Object.keys(line)[1] }}</div>
          <div class="section_field_value">
            {{ toHex(Object.values(line)[1]) }}
          </div>
          <div class="section_field_description">
            {{ getDescription(Object.keys(line)[1], Object.values(line)[1]) }}
          </div>
        </div>
      </div>
    </div>
    <div class="section">
      <div class="section_title">Data Directory</div>
      <div
        class="section_content"
        v-for="(dir, index) in Object.values(dataDirectory)"
        :key="index"
      >
        <div class="single">
          <div class="section_field_name">{{ getDirectoryName(index) }}</div>
          <div class="section_field_value">
            {{ toHex(dir["VirtualAddress"]) }}
          </div>
          <div class="section_field_description">
            {{ toHex(dir["Size"]) }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {
  dec2HexString,
  machine2String,
  unixtime2Human,
  fileCharacteristics2String,
  reverse,
  hex2a,
  dec2Hex,
} from "@/helpers/pe"
// import Copy from "@/components/elements/Copy"

export default {
  props: ["data"],
  components: {
    // copy: Copy,
  },
  computed: {
    fileHeader: function() {
      const data = this.data.FileHeader
      return [
        {
          type: "multiple",
          Machine: data.Machine,
          NumberOfSections: data.NumberOfSections,
        },
        {
          type: "multiple",
          Timestamp: data.TimeDateStamp,
          PointerToSymbolTable: data.PointerToSymbolTable,
        },
        {
          type: "multiple",
          NumberOfSymbols: data.NumberOfSymbols,
          SizeOfOptionalHeader: data.SizeOfOptionalHeader,
        },
        {
          type: "single",
          Characteristics: data.Characteristics,
        },
      ]
    },
    optionalHeader: function() {
      const data = this.data.OptionalHeader
      return [
        { type: "single", Magic: data.Magic },
        { type: "single", Subsystem: data.Subsystem },
        { type: "single", "Dll Characteristics": data.DllCharacteristics },
        { type: "single", CheckSum: data.CheckSum },
        { type: "single", "Address Of Entry Point": data.AddressOfEntryPoint },
        {
          type: "multiple",
          "Major Linker Version": data.MajorLinkerVersion,
          "Minor Linker Version": data.MinorLinkerVersion,
        },
        {
          type: "multiple",
          "Major Image Version": data.MajorImageVersion,
          "Minor Image Version": data.MinorImageVersion,
        },
        {
          type: "multiple",
          "Major Operating System Version": data.MajorOperatingSystemVersion,
          "Minor Operating System Version": data.MinorOperatingSystemVersion,
        },
        {
          type: "multiple",
          "Major Subsystem Version": data.MajorSubsystemVersion,
          "Minor Subsystem Version": data.MinorSubsystemVersion,
        },
        { type: "single", "Win32 Version Value": data.Win32VersionValue },
        {
          type: "multiple",
          ImageBase: data.ImageBase,
          "Size Of Image": data.SizeOfImage,
        },
        { type: "single", "Size Of Headers": data.SizeOfHeaders },
        {
          type: "multiple",
          "Base Of Code": data.BaseOfCode,
          "Size Of Code": data.SizeOfCode,
        },
        {
          type: "multiple",
          "Size Of Initialized Data": data.SizeOfInitializedData,
          "Size Of Uninitialized Data": data.SizeOfUninitializedData,
        },
        {
          type: "multiple",
          "Section Alignment": data.SectionAlignment,
          "File Alignment": data.FileAlignment,
        },
        {
          type: "multiple",
          "Size Of Stack Reserve": data.SizeOfStackReserve,
          "Size Of Stack Commit": data.SizeOfStackCommit,
        },
        {
          type: "multiple",
          "Size Of Heap Reserve": data.SizeOfHeapReserve,
          "Size Of Heap Commit": data.SizeOfHeapCommit,
        },
        { type: "single", LoaderFlags: data.LoaderFlags },
        { type: "single", "Number Of Rva And Sizes": data.NumberOfRvaAndSizes },
      ]
    },
    dataDirectory: function() {
      const data = this.data.OptionalHeader.DataDirectory
      return data
    },
    sections: function() {
      return {
        "File Header": this.fileHeader,
        "Optional Header": this.optionalHeader,
      }
    },
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
    getSize: function(value) {
      if (value >= 1000000) return (value / 1000000).toFixed(2) + " MB"
      if (value >= 1000) return (value / 1000).toFixed(2) + " KB"
      else return value + " B"
    },
    getDescription: function(key, value) {
      switch (key) {
        case "Machine":
          return machine2String(value)
        case "Timestamp":
          return unixtime2Human(value)
        case "Characteristics":
          return this._.join(fileCharacteristics2String(value), ", ")
        case "Magic":
          return reverse(hex2a(dec2Hex(value)))
        case "SizeOfOptionalHeader":
        case "Size Of Code":
        case "Size Of Initialized Data":
        case "Size Of Uninitialized Data":
        case "Section Alignment":
        case "File Alignment":
        case "Size Of Image":
        case "Size Of Headers":
        case "Size Of Stack Reserve":
        case "Size Of Stack Commit":
        case "Size Of Heap Reserve":
        case "Size Of Heap Commit":
          return this.getSize(value)
        default:
          return ""
      }
    },
    getDirectoryName: function(index) {
      switch (index) {
        case 0:
          return "Export Directory"
        case 1:
          return "Import Directory"
        case 2:
          return "Resource Directory"
        case 3:
          return "Exception Directory"
        case 4:
          return "Security Directory"
        case 5:
          return "Base Relocation Table"
        case 6:
          return "Debug Directory"
        case 7:
          return "Architecture specific"
        case 8:
          return "RVA of GlobalPointer"
        case 9:
          return "TLS Directory"
        case 10:
          return "Load Config Directory"
        case 11:
          return "Bound Import Directory"
        case 12:
          return "Import Address Table"
        case 13:
          return "Delay Import Descriptors"
        case 14:
          return "COM Runtime Descriptor"
        case 15:
          return "#15"
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.section {
  padding: 1rem;
  .section_title {
    font-size: large;
    color: #00d1b2;
  }
  .section_content {
    padding: 0.2rem;
    .multiple {
      display: flex;
      .first_field,
      .second_field {
        display: flex;
        .section_field_name {
          margin-left: 2rem;
          width: 15rem;
          font-weight: 500;
        }
        .section_field_value {
          width: 10rem;
          text-align: right;
        }
        .section_field_description {
          margin-left: 1rem;
          width: 20rem;
        }
      }
    }
    .single {
      display: flex;
      .section_field_name {
        margin-left: 2rem;
        width: 15rem;
        font-weight: 500;
      }
      .section_field_value {
        width: 10rem;
        text-align: right;
      }
      .section_field_description {
        margin-left: 1rem;
        width: 20rem;
      }
    }
  }
}
</style>
