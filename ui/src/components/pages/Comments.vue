<template>
  <div>
    <p id="noComments" v-if="usersData.length === 0">
      No Comments Available
    </p>
    <div v-for="(comment, index) in comments" :key="comment.id">
      <CommentCard
        :data="getCommentData(index)"
        v-if="getCommentData(index).avatar"
      />
    </div>
    <TextEditor v-if="this.$store.getters.getLoggedIn" />
    <div v-else id="authToCommentMessage" class="box">
      <span>Want to leave a comment? </span>
      <span>
        <router-link
          :to="{
            name: 'login',
            params: {
              nextUrl: $route.path,
            },
          }"
          >Sign In</router-link
        >
        &nbsp; or &nbsp;
        <router-link
          :to="{
            name: 'signUp',
            params: {
              nextUrl: $route.path,
            },
          }"
          >Sign Up</router-link
        >
      </span>
    </div>
  </div>
</template>

<script>
import TextEditor from "../elements/comments/TextEditor"
import CommentCard from "../elements/comments/CommentCard"

import { mapGetters } from "vuex"

export default {
  components: {
    TextEditor,
    CommentCard,
  },
  data() {
    return {
      usersData: [],
    }
  },
  watch: {
    comments: function() {
      this.loadUsers()
    },
  },
  computed: {
    ...mapGetters({
      comments: "getComments",
    }),
  },
  methods: {
    getUsersList: function() {
      return this._.uniq(this.comments.map((comment) => comment.username))
    },
    getUserData: async function(username) {
      if (username === this.$store.getters.getUsername) {
        var userData = this.$store.getters.getUserData
        var data = {
          username: username,
          name: userData.name,
          location: userData.location,
          avatar: this.$store.getters.getAvatar,
        }
        this.usersData.push(data)
        return
      }
      this.$http
        .get(this.$api_endpoints.USERS + username)
        .then((res) => {
          this.$http
            .get(this.$api_endpoints.USERS + username + "/avatar", {
              responseType: "arraybuffer",
            })
            .then((secRes) => {
              var data = {
                username: res.data.username,
                name: res.data.name,
                location: res.data.location,
                avatar: Buffer.from(secRes.data, "binary").toString("base64"),
              }
              this.usersData.push(data)
            })
        })
        .catch(() => {
          this.$awn.alert("An Error Occured While fetshing the user data")
        })
    },
    loadUsers: function() {
      var users = this.getUsersList()
      for (var index in users) {
        this.getUserData(users[index])
      }
    },
    getCommentData: function(index) {
      var comment = this.comments[index]
      var user = this.usersData.find(
        (user) => user.username === comment.username,
      )
      return this._.merge(comment, user)
    },
  },
  mounted() {
    if (this.comments && this.comments.length > 0) {
      this.loadUsers()
    }
  },
}
</script>

<style lang="scss">
#noComments {
  margin-bottom: 2em;
  font-size: 1.3em;
}
.comment_body {
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
}
#authToCommentMessage {
  margin-top: 2em;
  text-align: center;
  vertical-align: middle;
  display: grid;
  padding: 1em;
  span {
    font-size: 1.5em;
  }
}
</style>
