<template>
  <div>
    <div class="columns" v-if="!verification">
      <div class="column is-1 left_column">
        <img :src="'data:image/png;base64,' + avatar" />
        <div class="username">{{ this.data.username }}</div>
        <div>{{ this.time }}</div>
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
    <hr />
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
    },
    formatTimestamp: function() {
      var date = this._.replace(this.data.timestamp.substring(0, 10), ":", "")
      var time = this._.replace(this.data.timestamp.substring(11, 19), ":", "")
      this.time = moment(date + time, "YYYYMMDDhhmmss").fromNow()
    },
    deleteComment: function() {
      this.$http.delete(this.$api_endpoints.FILES+this.$store.getters.getHashContext+"/comments/"+this.data.id)
      .then(()=>{
          this.$store.dispatch('updateComments')
      })
      .catch(()=>{
          this.$awn.alert('An Error Occured While Deleting the comment, try again')
      })
    },
  },
  mounted() {
    if (this.data.username !== this.$store.getters.getUsername) {
      this.getAvatar()
      this.deletable = false
    } else {
      this.avatar = this.$store.getters.getAvatar
      this.deletable = true
    }
    this.formatTimestamp()
  },
}
</script>

<style lang="scss" scoped>
img {
  width: 5em;
}
h1 {
  font-size: 2rem;
}
h2 {
  font-size: 1.5rem;
}
blockquote {
  margin-left: 32px;
  border-left: 4px solid #ccc;
  padding-left: 8px;
}
.ql-syntax {
  background-color: #23241f;
  color: #f8f8f2;
  overflow: visible;
  white-space: pre-wrap;
  margin-bottom: 5px;
  margin-top: 5px;
  padding: 5px 10px;
}
ol {
  padding-left: 1.5em;
}
ul {
  padding-left: 1.5em;
}
hr {
  background-color: #dde0e3;
}
.username {
  font-size: 1.3em;
}
.left_column {
  text-align: center;
}
.right_column {
  .comment_body {
    vertical-align: center;
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
  button {
    font-size: 1.5em;
  }
  .danger {
    color: red;
    &:hover {
      color: red!important;
    }
  }
}
</style>
