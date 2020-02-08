<template>
  <button class="button is-outlined is-link" :class="{'active': liked}">
    <div class="rescan-button-text">
      <span class="icon">
        <i :class="{'fas fa-thumbs-up': !liked, 'fas fa-thumbs-down':liked}"></i>
      </span>
      <span>
        {{this.Liked?"Unlike":"Like"}}
      </span>
    </div>
  </button>
</template>

<script>
import { mapGetters } from "vuex"

export default {
  props: ["hash"],
  computed: {
    ...mapGetters({ userData: "getUserData" }),
    liked: function(){
        if(!this.userData.likes) return false
        return this.userData.likes.includes(this.hash)
    }
  },
  methods: {
    likeUnlike: function() {},
  },
}
</script>

<style></style>
