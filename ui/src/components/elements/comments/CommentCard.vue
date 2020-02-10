<template>
  <div class="comment box">
    <div class="columns" v-if="!verification">
      <div class="column is-2 left_column">
        <img :src="'data:image/png;base64,' + avatar" />
        <div class="username">{{ this.data.username }}</div>
        <div class="info">
          <div>
            <i class="icon fas fa-location-arrow"></i>
            {{ this.userData.location }}
          </div>
          <div>
            <i class="icon fas fa-clock"></i>
            {{ this.time }}
          </div>
        </div>
      </div>
      <div class="column is-1 separator">
        <hr />
      </div>
      <div class="column right_column">
        <div class="comment_body" v-html="data.body"></div>
      </div>
      <a class="delete" v-if="deletable" @click="verification = true"></a>
    </div>
    <div class="column verif" v-if="verification">
      <div>Are you Sure?</div>
      <button class="button is-light danger" @click="deleteComment">Yes</button>
      <button class="button is-light" @click="verification = false">
        No
      </button>
    </div>
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
      this.time = moment(this.data.timestamp).fromNow()
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
        .get(this.$api_endpoints.USERS + this.data.username, {
          responseType: "arraybuffer",
        })
        .then((res) => {
          this.userData = res
        })
        .catch(() => {
          this.$awn.alert("An Error Occured While fetshing the user data")
        })
    },
  },
  mounted() {
    if (this.data.username !== this.$store.getters.getUsername) {
      this.getAvatar()
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
img {
  width: 4em;
}

.username {
  font-size: 1.3em;
}

.left_column {
  text-align: center;
}

.right_column {
  padding: 15px;
}

.separator {
  width: fit-content !important;
  hr {
    background-color: #4a4a4a54;
    width: 1px;
    height: 100%;
    margin: auto;
  }
}

.columns {
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

.info {
  width: fit-content;
  margin: auto;
  text-align: left;
}
</style>
