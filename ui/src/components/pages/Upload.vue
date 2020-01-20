<template>
  <div class="columns" style="margin-top:100px">
    <div class="column is-8 is-offset-2">
      <tabs @tabChanged="tabChanged" type="is-boxed">
        <tab name="File" icon="ion-android-folder-open" :selected="true">
          <DropZone @fileAdded="onFileAdded" :enabled="ongoingStep === 0" />
        </tab>
        <!--<tab name="Url" icon="ion-ios-world-outline" >
          <form class="tile is-child box">
            <h1>This part is about our culture</h1>

            <p class="is-centered" style="margin-top:10px;">
              <small>
                By using Saferwall you consent to our
                <router-link :to="this.$routes.HOME.path"
                  >Terms of Service</router-link
                >
                and <router-link to="">Privacy Policy</router-link> and allow us
                to share your submission with the security community.
                <router-link to="">Learn more.</router-link>
              </small>
            </p>
          </form>
        </tab>-->
      </tabs>
      <progress-tracker alignment="center" v-if="selectedTab != 'Url'">
        <!-- prettier-ignore -->
        <step-item
          :title="stepTitle"
          v-for="(stepTitle, step) in {1: 'Uploaded', 2: 'Queued', 3: 'Processing', 4: 'Finished'}"
          :is-complete="Number(step) < ongoingStep"
          :is-active="Number(step) === ongoingStep"
          :key="step"
        ></step-item>
      </progress-tracker>
    </div>
  </div>
</template>
<script>
import Tabs from "@/components/elements/Tabs"
import Tab from "@/components/elements/Tab"
import DropZone from "@/components/elements/DropZone"
import ProgressTracker, { StepItem } from "vue-bulma-progress-tracker"

const step = {
  UPLOADED: 1,
  QUEUED: 2,
  PROCESSING: 3,
  FINISHED: 4,
  READY: 5,
}

export default {
  data() {
    return {
      selectedTab: "File",
      uploading: false,
      ongoingStep: 0, // by default, no step has started yet (0), next we move to step 1, 2 and so on
      pollInterval: null,
    }
  },
  components: {
    tabs: Tabs,
    tab: Tab,
    DropZone,
    ProgressTracker,
    StepItem,
  },
  methods: {
    tabChanged(selectedTab) {
      this.selectedTab = selectedTab
    },
    onFileAdded(file) {
      if (!file) {
        this.$awn.alert("File cannot be read!")
        return
      }

      // check if size exceeds 64mb
      if (file.size > 64000000) {
        this.$awn.alert("File size exceeds 64MB !")
        return
      }
      const reader = new FileReader()
      reader.onload = (loadEvent) => {
        // file has been read successfully
        const fileBuffer = loadEvent.target.result
        crypto.subtle
          .digest("SHA-256", fileBuffer)
          .then((hashBuffer) => {
            const hashArray = new Uint8Array(hashBuffer)
            let hashHex = ""
            for (let i = 0; i < hashArray.byteLength; i++) {
              let hex = new Number(hashArray[i]).toString("16")
              if (hex.length === 1) {
                hex = "0" + hex
              }
              hashHex += hex
            }
            // hash hexadecimal has been calculated successfully
            this.$http
              .get(`${this.$api_endpoints.FILES}${hashHex}/`)
              .then((response) => {
                // file exists
                this.$router.push(this.$routes.SUMMARY.path + hashHex)
              })
              .catch(() => {
                this.ongoingStep = step.UPLOADED
                const formData = new FormData()
                formData.append("file", file)
                this.$http
                  .post(this.$api_endpoints.FILES, formData, {
                    headers: {
                      "Content-Type": "multipart/form-data",
                    },
                  })
                  .then(() => {
                    // set a poll interval of 5s

                    this.pollInterval = setInterval(
                      this.fetchStatus,
                      3000,
                      hashHex,
                    )
                    setTimeout(() => {
                      clearInterval(this.pollInterval)
                      this.trackException()
                    }, 36000000) // stop polling after an hour
                  })
                  .catch(console.log)
              })
          })
          .catch(() => {
            this.$awn.alert(
              "Sorry, we couldn't upload the file. Please, try again!",
            )
            this.trackException()
          })
      }
      reader.readAsArrayBuffer(file)
    },
    fetchStatus(hashHex) {
      this.$http
        .get(`${this.$api_endpoints.FILES}${hashHex}/`)
        .then((response) => {
          const status = response.data.status
          // change ongoingStep according to status
          // status
          // 0: queued
          // 1: processing
          // 2: finished
          switch (status) {
            case 0:
              this.ongoingStep = step.ENQUEUE
              break
            case 1:
              this.ongoingStep = step.PROCESSING
              break
            case 2:
              this.ongoingStep = step.FINISHED
              // stop polling
              clearInterval(this.pollInterval)
              setTimeout(() => {
                this.ongoingStep = step.READY
                this.$router.push({
                  name: this.$routes.SUMMARY.name,
                  params: { hash: hashHex },
                })
              }, 4000)
              this.trackSuccess()
              this.$store.dispatch("updateHash", hashHex)
              this.$store.dispatch("updateFileData", response)
              this.$router.push(this.$routes.SUMMARY.path)
              break
          }
        })
        .catch(() => {
          this.$awn.alert("Problem occured while uploading, try again")
          this.trackException()
        })
    },
    trackSuccess() {
      this.$gtag.event("Upload_Success", {
        event_category: "Upload",
        event_label: "Success Upload",
        value: 1,
      })
    },
    trackException() {
      this.$gtag.exception({
        description: "Failed Upload ",
        fatal: false,
      })
    },
  },
}
</script>
<style lang="scss" scoped>
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}
.slide-enter,
.slide-leave-to {
  height: 0;
  padding: 0;
  overflow: hidden;
  opacity: 0;
}

.progress {
  margin-top: 1em;
}

.tile {
  margin-top: -24px !important;
  border-radius: 0 0 4px 4px !important;
  margin-left: 1px !important;
}

a.button {
  display: inline-block;
  margin-top: 10px;
}
</style>
