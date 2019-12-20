<template>
  <div>
    <button class="button is-outlined" @click="downloadFile">
      <span class="icon">
        <i class="fas fa-file-download"></i>
      </span>
      <span>
        Download File
      </span>
    </button>
  </div>
</template>

<script>
export default {
  methods: {
    downloadFile: function() {
      this.$http
        .get(
          `${this.$api_endpoints.FILES}${this.$store.getters.getHashContext}/download/`,
          {
            responseType: "blob",
          },
        )
        .then((response) => {
          console.log(response)
          const url = window.URL.createObjectURL(new Blob([response.data]))
          const link = document.createElement("a")
          link.href = url
          link.setAttribute("download", `${this.$store.getters.getHashContext}.zip`)
          document.body.appendChild(link)
          link.click()
        })
        .catch((e) => console.log(e))
    },
  },
}
</script>

<style></style>
