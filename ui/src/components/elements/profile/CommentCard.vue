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
            {{ this.liked ? "Unlike" : "Like" }}
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
      this.$store.dispatch("updateHash", this.comment.sha256)
      this.$router.push(this.$routes.SUMMARY.path + this.comment.sha256)
    },
  },
  mounted() {
    if (this.$store.getters.getLikes.includes(this.comment.sha256))
      this.liked = true
    this.time = moment(this.comment.timestamp).format("MMMM Do YYYY")
  },
}
</script>

<style lang="scss">
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
        background-color: rgba(0, 0, 0, 0.03);
        border-radius: 10px;
        padding: 1em;
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
