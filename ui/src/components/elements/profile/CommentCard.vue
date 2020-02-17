<template>
  <div class="level">
    <div class="level-left">
      <div class="level-item">
        <div id="hash" @click="showFile">{{ comment.sha256 }}</div>
      </div>
      <div class="level-item">
        <span id="comment_body" v-html="comment.body"></span>
      </div>
      <div class="level-item">
        <span id="timestamp">
          <i class="icon fas fa-clock"></i>
          {{ time }}
        </span>
      </div>
    </div>
    <div class="level-right">
      <div class="level-item">
        <button class="button" :class="{ active: liked }" @click="likeUnlike">
          <span class="icon">
            <i class="fas fa-heart"></i>
          </span>
          <span>
            {{ this.liked ? "unlike" : "like" }}
          </span>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import moment from "moment"

export default {
  props: ["comment"],
  data() {
    return {
      liked: false,
      time: null,
    }
  },
  methods: {
    likeUnlike: function() {
      this.$http
        .post(`${this.$api_endpoints.FILES}${this.comment.sha256}/actions/`, {
          type: this.liked ? "unlike" : "like",
        })
        .then(() => {
          this.liked = !this.liked
          this.$store.dispatch("updateLikes")
        })
        .catch(() => {
          this.$awn.alert("An Error Occured, try again")
        })
    },
    showFile: function() {
      this.$store.dispatch("updateHash", this.file.sha256)
      this.$router.push(this.$routes.SUMMARY.path + this.file.sha256)
    },
  },
  mounted() {
    if (this.$store.getters.getLikes.includes(this.comment.sha256))
      this.liked = true
    this.time = moment(this.comment.timestamp).format("MMMM Do YYYY")
  },
}
</script>

<style lang="scss" scoped>
.level {
  padding: 1em;
  border-bottom-color: #dbdbdb;
  border-bottom-style: solid;
  border-bottom-width: 1px;
  .level-left {
    display: block;
    .level-item {
      justify-content: left;
      #comment_body {
        padding: 0.3em;
      }
      #hash {
        font-size: large;
        font-weight: 500;
        cursor: pointer;
      }
      svg {
        vertical-align: bottom;
      }
    }
  }
  .button {
    background-color: transparent;
    border-color: #f14668;
    span {
      color: #f14668;
    }

    &:hover {
      background-color: #f14668;
      span {
        color: white;
      }
    }

    &.active {
      background-color: #f14668;
      span {
        color: white;
      }
    }
  }
}
</style>
