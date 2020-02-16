<template>
  <div class="media box" @click="goToProfile">
    <figure class="media-left">
      <p class="image is-64x64">
        <img id="avatar" :src="'data:image/png;base64,' + userData.avatar" />
      </p>
    </figure>
    <div class="media-content">
      <div class="content">
        <p>
          <strong>{{
            this.userData.name ? this.userData.name : this.userData.username
          }}</strong>
          <small> @{{ this.userData.username }}</small>
          <br />
          {{ this.userData.location }}
        </p>
      </div>
    </div>
    <div class="media-right">
      <button :disabled="self" class="button">{{this.followed?"Unfollow":"Follow"}}</button>
    </div>
  </div>
</template>

<script>
export default {
  props: ["userData"],
  data() {
    return {
      followed: false,
      self: false,
    }
  },
  methods: {
    followUnfollow: function() {
      if (this.followed) this.unfollow()
      else this.follow()
    },
    follow: function() {
      this.$http
        .post(
          this.$api_endpoints.USERS + this.userData.username + "/actions/",
          {
            type: "follow",
          },
        )
        .then(() => {
          this.followed = true
        })
        .catch()
    },
    unfollow: function() {
      this.$http
        .post(
          this.$api_endpoints.USERS + this.userData.username + "/actions/",
          {
            type: "unfollow",
          },
        )
        .then(() => {
          this.followed = false
        })
        .catch()
    },
    goToProfile: function() {
      this.$router.replace({ name: 'profile', params: { user: this.userData.username } })
    },
  },
  mounted() {
    if (this.$store.getters.getFollowing.includes(this.userData.username))
      this.followed = true
    if (this.userData.username === this.$store.getters.getUsername)
      this.self = true
  },
}
</script>

<style lang="scss" scoped>
.media {
  cursor: pointer;
  .media-left {
    justify-content: left;
  }
}
</style>
