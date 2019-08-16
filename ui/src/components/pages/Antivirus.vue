<template>
  <div class="tile is-ancestor">
    <loader v-if="showLoader"></loader>
    <div class="tile is-parent is-6">
      <div class="tile is-child box" v-if="!showLoader">
        <h4 class="title">First Scan</h4>
        <table class="table is-striped is-fullwidth">
          <thead>
            <tr>
              <th>Antivirus</th>
              <th>Output</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(value, vendor) in firstScan" :key="vendor">
              <td>{{ vendor }}</td>
              <td>
                <span
                  :class="[
                    { 'has-text-success': !value.detected },
                    { 'has-text-danger': value.detected }
                  ]"
                  style="position:relative"
                  @mouseover="mouseOver('first', vendor)"
                  @mouseleave="mouseLeave('first', vendor)"
                >
                  <span
                    :class="{
                      transparent:
                        value.detected &&
                        JSON.stringify(show) ===
                          JSON.stringify({ type: 'first', vendor: vendor })
                    }"
                  >
                    <i
                      class="output-icon icon"
                      :class="[
                        { 'ion-alert-circled': value.detected },
                        { 'ion-checkmark-circled': !value.detected }
                      ]"
                    ></i>
                    {{ value.result || "Clean" }}
                  </span>
                  <transition name="fade">
                    <copy
                      v-if="
                        value.detected &&
                          JSON.stringify(show) ===
                            JSON.stringify({ type: 'first', vendor: vendor })
                      "
                      :content="value.result"
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
      <div class="tile is-child box" v-if="!showLoader">
        <h4 class="title">Last Scan</h4>
        <table class="table is-striped is-fullwidth">
          <thead>
            <tr>
              <th>Antivirus</th>
              <th>Output</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(value, vendor) in lastScan" :key="vendor">
              <td>{{ vendor }}</td>
              <td>
                <span
                  :class="[
                    { 'has-text-success': !value.detected },
                    { 'has-text-danger': value.detected }
                  ]"
                  style="position:relative"
                  @mouseover="mouseOver('last', vendor)"
                  @mouseleave="mouseLeave('last', vendor)"
                >
                  <span
                    :class="{
                      transparent:
                        value.detected &&
                        JSON.stringify(show) ===
                          JSON.stringify({ type: 'last', vendor: vendor })
                    }"
                  >
                    <i
                      class="output-icon icon"
                      :class="[
                        { 'ion-alert-circled': value.detected },
                        { 'ion-checkmark-circled': !value.detected }
                      ]"
                    ></i>
                    {{ value.result || "Clean" }}
                  </span>
                  <transition name="fade">
                    <copy
                      v-if="
                        value.detected &&
                          JSON.stringify(show) ===
                            JSON.stringify({ type: 'last', vendor: vendor })
                      "
                      :content="value.result"
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
import axios from 'axios'
import Loader from '@/components/elements/Loader'
import Copy from '@/components/elements/Copy'

export default {
  components: {
    loader: Loader,
    copy: Copy
  },
  data () {
    return {
      showLoader: true,
      lastScan: {},
      firstScan: {},
      show: { type: '', vendor: '' }
    }
  },
  methods: {
    mouseOver (type, vendor) {
      this.show = { type, vendor }
    },
    mouseLeave (type, index) {
      this.show = {}
    }
  },
  mounted () {
    // replace route params with props
    axios
      .get(`/api/v1/files/${this.$route.params.hash}`)
      .then(data => {
        this.showLoader = false
        if (!data.data.multiav) {
          return
        }
        this.firstScan = data.data.multiav.first_scan
        this.lastScan = data.data.multiav.last_scan

        Object.keys(this.firstScan).forEach(key => {
          const first = this.firstScan[key]
          first.showCopy = false
          const last = this.lastScan[key]
          last.showCopy = false
        })
      })
      .catch(err => console.error(err))
  }
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
