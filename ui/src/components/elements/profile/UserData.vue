<template>
  <div class="tile is-vertical">
    <img
      id="Profile_avatar"
      :src="'data:image/png;base64,' + userData.avatar"
    />
    <div id="name">
      {{ this.userData.name ? this.userData.name : this.userData.username }}
    </div>
    <div id="username">@{{ this.userData.username }}</div>
    <div id="bio" v-if="this.userData.bio">{{ this.userData.bio }}</div>
    <div id="location" v-if="this.userData.location">
      <i class="icon fas fa-location-arrow"></i>
      {{ this.userData.location }}
    </div>
    <div id="url" v-if="this.userData.url">
      <i class="icon fas fa-link"></i>
      {{ this.userData.url }}
    </div>
    <button
      class="button is-medium is-primary"
      id="follow"
      @click="followUnfollow"
      :disabled="self"
    >
      {{ this.followed ? "Unfollow" : "Follow" }}
    </button>
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
  watch: {
    userData: function() {
      if (this.$store.getters.getFollowing.includes(this.userData.username))
        this.followed = true
      if (this.userData.username === this.$store.getters.getUsername)
        this.self = true
    },
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
.tile {
  align-items: center;
  padding: 0.5em;
  #Profile_avatar {
    width: 50%;
  }
  #name {
    font-size: x-large;
    font-weight: bold;
  }
  #username {
    font-size: large;
    font-weight: lighter;
  }
  #follow {
    width: 50%;
  }
  #bio {
    font-size: initial;
  }
  #location {
    align-items: center;
  }
  .button{
    margin-top: 0.5em;
  }
}
</style>
