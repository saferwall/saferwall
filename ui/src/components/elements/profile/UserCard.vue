<template>
  <div class="media">
    <figure class="media-left">
      <p class="image is-64x64">
        <img id="avatar" :src="'data:image/png;base64,' + userData.avatar" />
      </p>
    </figure>
    <div class="media-content">
      <div class="content">
        <p>
          <strong @click="goToProfile" id="username">{{
            this.userData.name ? this.userData.name : this.userData.username
          }}</strong>
          <small> @{{ this.userData.username }}</small>
          <br />
          {{ this.userData.location }}
        </p>
      </div>
    </div>
    <div class="media-right">
      <button :disabled="self" class="button" @click="followUnfollow">
        {{ this.followed ? "Unfollow" : "Follow" }}
      </button>
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
      if (!this.$store.getters.getLoggedIn) {
        this.$router.push({
          name: "login",
          params: {
            nextUrl: this.$route.path,
          },
        })
        return
      }
      this.$http
        .post(
          this.$api_endpoints.USERS + this.userData.username + "/actions/",
          {
            type: this.followed ? "unfollow" : "follow",
          },
        )
        .then(() => {
          this.followed = !this.followed
          this.$store.dispatch("updateFollowing")
        })
        .catch()
    },
    goToProfile: function() {
      this.$router.replace({
        name: "profile",
        params: { user: this.userData.username },
      })
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
  padding: 1em;
  border-bottom-color: #dbdbdb;
  border-bottom-style: solid;
  border-bottom-width: 1px;
  #username {
    cursor: pointer;
  }
  .media-left {
    justify-content: left;
  }
}
</style>
