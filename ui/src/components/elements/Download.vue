<template>
    <button class="button is-info" @click="downloadFile">
      <span class="icon">
        <i class="fas fa-file-download"></i>
      </span>
      <span>
        Download File
      </span>
    </button>
</template>

<script>
export default {
  props:['hash'],
  methods: {
    downloadFile: function() {
      this.$http
        .get(
          `${this.$api_endpoints.FILES}${this.hash}/download/`,
          {
            responseType: "blob",
          },
        )
        .then((response) => {
          const url = window.URL.createObjectURL(new Blob([response.data]))
          const link = document.createElement("a")
          link.href = url
          link.setAttribute(
            "download",
            `${this.$store.getters.getHashContext}.zip`,
          )
          document.body.appendChild(link)
          link.click()
          this.track()
        })
        .catch((e) => console.log(e))
    },
    track() {
      this.$gtag.event("Download_Success", {
        event_category: "Download",
        event_label: "File Downloaded, hash: "+this.hash,
        value: 1,
      })
    },
  },
}
</script>

<style></style>
