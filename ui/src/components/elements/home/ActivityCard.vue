<template>
  <div class="media">
    <figure class="media-left">
      <p class="image is-64x64">
        <img id="avatar" :src="this.userData ?'data:image/png;base64,'+this.userData.avatar:''" />
      </p>
    </figure>
    <div class="media-content">
      <div class="content">
        <p>
          <strong @click="goToProfile" id="username">{{ username }}</strong>
          &nbsp;
          <span class="action">{{ this.action }}</span>
          &nbsp;
          <strong v-if="activity.type === 'follow'">{{
            this.activity.content?this.activity.content.user:''
          }}</strong>
          &nbsp;
          <small id="time">{{ this.time }}</small>
          <br />
        </p>
      </div>
      <div v-if="activity.content && userData">
        <div class="content" v-if="this.activity.type === 'submit'">
          <FileCard :hash="activity.content.sha256" />
        </div>
        <div class="content" v-if="this.activity.type === 'like'">
          <FileCard :hash="activity.content.sha256" />
        </div>
        <div class="content" v-if="this.activity.type === 'follow'">
          <UserCard :username="activity.content.user" />
        </div>
        <div class="content" v-if="this.activity.type === 'comment'">
          <CommentCard :data="activity.content" />
        </div>
      </div>
    </div>
    <div class="media-right" v-if="this.$store.getLoggedIn">
      <button class="button" @click="followUnfollow">
        {{ this.followed ? "Unfollow" : "Follow" }}
      </button>
    </div>
  </div>
</template>

<script>
import moment from "moment"

import UserCard from "./UserCard"
import FileCard from "./FileCard"
import CommentCard from "./CommentCard"

export default {
  props: ["activity", "userData"],
  components: {
    UserCard,
    FileCard,
    CommentCard,
  },
  data() {
    return {
      action: "",
      followed: false,
      time: "",
      fileData: {},
    }
  },
  watch: {
    userData: function() {
      if (this.$store.getters.getFollowing.includes(this.userData.username))
        this.followed = true
    },
  },
  computed: {
    username: function() {
      return this.userData ? this.userData.username : ""
    },
  },
  methods: {
    setAction() {
      switch (this.activity.type) {
        case "comment":
          this.action = "Commented a file"
          break
        case "like":
          this.action = "Liked a file"
          break
        case "submit":
          this.action = "Submitted a file"
          break
        case "follow":
          this.action = "Started following"
          break
      }
    },
    formatTime() {
      this.time = moment(this.activity.timestamp).fromNow()
    },
    followUnfollow: function() {
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
    this.setAction()
    this.formatTime()
  },
}
</script>

<style lang="scss" scoped>
#username {
  cursor: pointer;
}
#time {
  padding-left: 0.3em;
  font-weight: lighter;
}
#avatar {
  margin: auto;
  width: 70%;
}
</style>
