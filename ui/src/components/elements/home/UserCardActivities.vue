<template>
  <div class="level tile">
    <div class="image is-64x64">
      <img
        id="avatar"
        :src="userData.avatar ? 'data:image/png;base64,' + userData.avatar : ''"
      />
    </div>
    <div class="media-content">
      <div class="content">
        <strong @click="goToProfile" id="username">{{
          this.userData.username
        }}</strong>
        &nbsp;
        <small>@{{ this.userData.username }}</small>
      </div>
    </div>
    <div class="btn">
      <button
        class="button"
        v-if="this.$store.getters.getLoggedIn"
        @click="followUnfollow"
        :disabled="self"
      >
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
      this.$http
        .post(this.$api_endpoints.USERS + this.userData.username + "/actions/", {
          type: this.followed ? "unfollow" : "follow",
        })
        .then(() => {
          this.followed = !this.followed
        })
        .catch()
    },
    goToProfile: function() {
      this.$router.push({
        name: "profile",
        params: { user: this.userData.username },
      })
    },
  },
  mounted() {
    if (this.$store.getters.getFollowing.includes(this.userData.username))
      this.followed = true
    if (this.username === this.$store.getters.getUsername) this.self = true
  },
}
</script>

<style lang="scss" scoped>
.tile {
  border: 1px solid !important;
  border-color: #d1d5da !important;
  margin-left: 2em;
  margin-right: 2em;
  width: 90%;
  padding: 0.5em;
  align-items: center;

  .content {
    display: table;
    padding-right: 1.5em;
    & * {
      display: table-row;
    }
  }

  #username {
    cursor: pointer;
  }

  & .btn {
    right: 0%;
  }

  .image {
    display: flex;
    align-items: center;
    #avatar {
      width: 70%;
      margin: auto;
    }
  }

  .button {
    border-color: #d1d5da;

    &:hover {
      background-color: #d1d5da;
    }

    &[disabled] {
      &:hover {
        background-color: transparent;
      }
    }

    &.active {
      background-color: #d1d5da;
    }
  }
}
</style>
