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
          this.userData.name ? this.userData.name : this.username
        }}</strong>
        &nbsp;
        <small>@{{ this.username }}</small>
      </div>
    </div>
    <div class="btn">
      <button
        class="button"
        v-if="this.$store.getters.getLoggedIn"
        @click="followUnfollow"
      >
        {{ this.followed ? "Unfollow" : "Follow" }}
      </button>
    </div>
  </div>
</template>

<script>
export default {
  props: ["username"],
  data() {
    return {
      userData: {},
      followed: false,
    }
  },
  methods: {
    getUserData() {
      this.$http
        .get(this.$api_endpoints.USERS + this.username + "?fields=name")
        .then((res) => {
          this.$http
            .get(this.$api_endpoints.USERS + this.username + "/avatar", {
              responseType: "arraybuffer",
            })
            .then((secRes) => {
              var data = {
                name: res.data.name,
                avatar: Buffer.from(secRes.data, "binary").toString("base64"),
              }
              this.userData = data
            })
        })
        .catch(() => {
          this.$awn.alert("An Error Occured While fetshing the user data")
        })
    },
    followUnfollow: function() {
      this.$http
        .post(this.$api_endpoints.USERS + this.username + "/actions/", {
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
        params: { user: this.username },
      })
    },
  },
  mounted() {
    this.getUserData()
    if (this.$store.getters.getFollowing.includes(this.userData.username))
      this.followed = true
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

    &.active {
      background-color: #d1d5da;
    }
  }
}
</style>
