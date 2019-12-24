<template>
  <div v-if="this.route !== 'upload'">
    <button
      class="button is-outlined"
      @click="rescanFile"
      :disabled="Rescanning"
      is-loading
    >
      <div v-if="!Rescanning">
        <span class="icon">
          <i class="fas fa-redo-alt"></i>
        </span>
        <span>
          Rescan File
        </span>
      </div>
      <div v-if="Rescanning">
        <span>
          {{ this.stepText }}
        </span>
      </div>
    </button>
  </div>
</template>

<script>
const step = {
  ENQUEUE: 0,
  PROCESSING: 2,
  FINISHED: 3,
}

export default {
  props: ["route"],
  data() {
    return {
      hash: "",
      ongoingStep: 0,
      Rescanning: false,
      pollInterval: null,
    }
  },
  computed: {
    stepText: function() {
      switch (this.ongoingStep) {
        case 0:
          return "Queud"
        case 2:
          return "Processing"
        case 3:
          return "Finished"
        default:
          return ""
      }
    },
  },
  methods: {
    rescanFile: function() {
      this.$http
        .post(`${this.$api_endpoints.FILES}${this.hash}/actions/`, {
          type: "rescan",
        })
        .then(() => {
          this.$awn.async(this.setPollInterval())
        })
        .catch()
    },
    setPollInterval() {
      this.Rescanning = true
      this.pollInterval = setInterval(this.fetchStatus, 3000)
      return new Promise((resolve) => {
        setTimeout(() => {
          this.$awn.alert("Something went wrong, try again")
          clearInterval(this.pollInterval)
        }, 36000000) // stop polling after an hour
      })
    },
    fetchStatus() {
      this.$http
        .get(`${this.$api_endpoints.FILES}${this.hash}/`)
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
                  name: this.$router.SUMMARY.name,
                  params: { hash: this.hash },
                })
              }, 4000)
              this.$store.dispatch("updateHash", this.hash)
              this.$store.dispatch("updateFileData", response)
              this.Rescanning = false
              this.ongoingStep = 0
              this.$awn.closeToasts()
              this.$awn.success("successfully rescaned the file")
              this.UpdateData()
              break
          }
        })
        .catch(() => {
          this.$awn.alert(
            "Problem occured while rescanning the file, try again",
          )
        })
    },
    UpdateData: function() {
      this.$http
        .get(this.$api_endpoints.FILES + this.hash)
        .then((data) => {
          this.$store.dispatch(
            "updateHash",
            this.$router.currentRoute.params.hash,
            this.$store.dispatch("updateFileData", data),
          )
          this.$awn.success("successfully updated the data")
        })
        .catch(() => {
          this.$awn.alert(
            "Error occured while updating file info, try refreshing the page",
          )
        })
    },
  },
  created() {
    this.hash = this.$store.getters.getHashContext
  },
}
</script>

<style></style>
