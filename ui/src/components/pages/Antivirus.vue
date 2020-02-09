<template>
  <div class="tile is-ancestor">
    <div class="tile is-parent is-6">
      <div class="tile is-child box">
        <h4 class="title">First Scan</h4>
        <table class="table is-striped is-fullwidth">
          <thead>
            <tr>
              <th>Antivirus</th>
              <th>Output</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(value, vendor) in Scans.firstScan" :key="vendor">
              <td>{{ vendor }}</td>
              <td>
                <span
                  :class="[
                    { 'has-text-success': !value.infected },
                    { 'has-text-danger': value.infected },
                  ]"
                  style="position:relative"
                  @mouseover="mouseOver('first', vendor)"
                  @mouseleave="mouseLeave('first', vendor)"
                >
                  <span
                    :class="{
                      transparent:
                        value.infected &&
                        JSON.stringify(show) ===
                          JSON.stringify({ type: 'first', vendor: vendor }),
                    }"
                  >
                    <i
                      class="output-icon icon"
                      :class="[
                        { 'ion-alert-circled': value.infected },
                        { 'ion-checkmark-circled': !value.infected },
                      ]"
                    ></i>
                    {{ value.output || "Clean" }}
                  </span>
                  <transition name="fade">
                    <copy
                      v-if="
                        value.infected &&
                          JSON.stringify(show) ===
                            JSON.stringify({ type: 'first', vendor: vendor })
                      "
                      :content="value.output"
                    ></copy>
                  </transition>
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <div class="tile is-parent is-6">
      <div class="tile is-child box">
        <h4 class="title">Last Scan</h4>
        <table class="table is-striped is-fullwidth">
          <thead>
            <tr>
              <th>Antivirus</th>
              <th>Output</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(value, vendor) in Scans.lastScan" :key="vendor">
              <td>{{ vendor }}</td>
              <td>
                <span
                  :class="[
                    { 'has-text-success': !value.infected },
                    { 'has-text-danger': value.infected },
                  ]"
                  style="position:relative"
                  @mouseover="mouseOver('last', vendor)"
                  @mouseleave="mouseLeave('last', vendor)"
                >
                  <span
                    :class="{
                      transparent:
                        value.infected &&
                        JSON.stringify(show) ===
                          JSON.stringify({ type: 'last', vendor: vendor }),
                    }"
                  >
                    <i
                      class="output-icon icon"
                      :class="[
                        { 'ion-alert-circled': value.infected },
                        { 'ion-checkmark-circled': !value.infected },
                      ]"
                    ></i>
                    {{ value.output || "Clean" }}
                  </span>
                  <transition name="fade">
                    <copy
                      v-if="
                        value.infected &&
                          JSON.stringify(show) ===
                            JSON.stringify({ type: 'last', vendor: vendor })
                      "
                      :content="value.output"
                    ></copy>
                  </transition>
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import Copy from "@/components/elements/Copy"
import { mapGetters } from "vuex"

export default {
  components: {
    copy: Copy,
  },
  data() {
    return {
      show: { type: "", vendor: "" },
    }
  },
  computed: {
    ...mapGetters({ fileData: "getFileData" }),
    Scans: function() {
      if (this.fileData === {} || !this.fileData) return {}

      var _firstScan = this.fileData.data.multiav.first_scan
      var _lastScan = this.fileData.data.multiav.last_scan

      Object.keys(_firstScan).forEach((key) => {
        _firstScan[key].showCopy = false
        _lastScan[key].showCopy = false
      })
      return {
        firstScan: _firstScan,
        lastScan: _lastScan,
      }
    },
  },
  methods: {
    mouseOver(type, vendor) {
      this.show = { type, vendor }
    },
    mouseLeave(type, index) {
      this.show = {}
    },
  },
}
</script>

<style lang="scss" scoped>
.fade-enter-active,
.fade-leave-active {
  transition-property: opacity;
  transition-duration: 0.25s;
}

.fade-enter-active {
  transition-delay: 0;
}

.fade-enter,
.fade-leave-active {
  opacity: 0;
}
span {
  transition: all 0.2s;
}
.transparent {
  opacity: 0.35;
}
.output-icon {
  font-size: 18px;
  display: inline-block;
}
</style>
