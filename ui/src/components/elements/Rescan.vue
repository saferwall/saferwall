<template>
  <div>
    <button class="button is-outlined" @click="rescanFile">
      <span class="icon">
        <i class="fas fa-redo-alt"></i>
      </span>
      <span>
        Rescan File
      </span>
    </button>
  </div>
</template>

<script>
export default {
  methods: {
    rescanFile: function() {
      this.$http
        .post(
          `${this.$api_endpoints.FILES}${this.$store.getters.getHashContext}/actions/`,
          {
            type: "rescan",
          },
        )
        .then((response) => {
          this.$router.push({name:"summary"})
        })
        .catch((e) => console.log(e))
    },
  },
}
</script>

<style></style>
