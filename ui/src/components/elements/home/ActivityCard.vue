<template>
  <div class="media">
      <figure class="media-left">
      <p class="image is-64x64">
        <img id="avatar" :src="'data:image/png;base64,'" />
      </p>
    </figure>
    <div class="media-content">
      <div class="content">
        
      </div>
    </div>
    <div class="media-right">
      <button class="button" >
        "Unfollow"
      </button>
    </div>
  </div>
</template>

<script>
export default {
  props: ["activity", "userData"],
  data() {
    return {}
  },
  methods: {},
  mounted() {},
}
</script>

<style></style>
