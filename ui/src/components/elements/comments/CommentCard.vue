<template>
  <div>
    <article class="media box" v-if="!verification">
      <figure class="media-left">
        <p class="image is-64x64">
          <img :src="'data:image/png;base64,' + avatar" />
        </p>
      </figure>
      <div class="media-content">
        <div class="content">
          <p>
            <strong class="username">{{ this.data.username }}</strong>
            &nbsp;
            <small>{{ this.time }}</small>
            <br />
            <span class="comment_body" v-html="data.body"></span>
          </p>
        </div>
      </div>
      <div class="media-right">
        <button
          class="delete"
          v-if="deletable"
          @click="verification = true"
        ></button>
      </div>
    </article>
    <article class="media box" v-if="verification">
      <div class="media-content">
        <div class="content verif">
          <div>Are you Sure?</div>
          <button class="button is-light danger" @click="deleteComment">
            Yes
          </button>
          <button class="button is-light" @click="verification = false">
            No
          </button>
        </div>
      </div>
    </article>
  </div>
</template>

<script>
import moment from "moment"

export default {
  props: ["data"],
  data() {
    return {
      avatar: null,
      time: null,
      deletable: false,
      verification: false,
      userData: {},
    }
  },
  methods: {
    getAvatar: function() {
      this.$http
        .get(this.$api_endpoints.USERS + this.data.username + "/avatar", {
          responseType: "arraybuffer",
        })
        .then((secRes) => {
          this.avatar = Buffer.from(secRes.data, "binary").toString("base64")
        })
        .catch(() => {
          this.$awn.alert("An Error Occured While fetshing the user avatar")
        })
    },
    formatTimestamp: function() {
      this.time = moment(this.data.timestamp).fromNow(true)
    },
    deleteComment: function() {
      this.$http
        .delete(
          this.$api_endpoints.FILES +
            this.$store.getters.getHashContext +
            "/comments/" +
            this.data.id,
        )
        .then(() => {
          this.$store.dispatch("updateComments")
        })
        .catch(() => {
          this.$awn.alert(
            "An Error Occured While Deleting the comment, try again",
          )
        })
    },
    getUserData: function() {
      this.$http
        .get(this.$api_endpoints.USERS + this.data.username)
        .then((res) => {
          this.userData = res.data
        })
        .catch(() => {
          this.$awn.alert("An Error Occured While fetshing the user data")
        })
    },
  },
  mounted() {
    if (this.data.username !== this.$store.getters.getUsername) {
      this.getAvatar()
      this.getUserData()
      this.deletable = false
    } else {
      this.avatar = this.$store.getters.getAvatar
      this.userData = this.$store.getters.getUserData
      this.deletable = true
    }
    this.formatTimestamp()
  },
}
</script>

<style scoped lang="scss">

.username {
  font-size: 1em;
  font-weight: 600;
}

.media {
  margin-bottom: 0.6em;

  .delete {
    right: 3%;
    position: absolute;
    display: none;
  }
  &:hover {
    .delete {
      display: block;
    }
  }
}

.verif {
  text-align: center;
  vertical-align: middle;
  div {
    font-size: 1.5em;
  }
  .button {
    font-size: 1.5em;
    background-color: transparent;
    &:hover {
      background-color: #00000012;
    }
  }
  .danger {
    color: red;
    &:hover {
      color: red !important;
    }
  }
}
</style>
