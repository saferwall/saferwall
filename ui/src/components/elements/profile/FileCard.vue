<template>
  <div class="level">
    <div class="level-left">
      <div class="level-item">
        <div id="hash" @click="showFile">{{ file.sha256 }}</div>
      </div>
      <div class="level-item">
        <span id="tags">
          <i class="icon fas fa-tags"></i>
          Tags:
          <span v-if="!file.tags">none</span>
          <span id="tag" v-for="tag in file.tags" :key="tag">{{ tag }}</span>
        </span>
        <span id="Av">
          <i class="icon fas fa-search"></i>
          Av Detection Count: {{ file.AvDetectionCount }}
        </span>
        <span id="timestamp" v-if="file.timestamp">
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
  props: ["file"],
  data() {
    return {
      liked: false,
      time: null,
    }
  },
  methods: {
    likeUnlike: function() {
      this.$http
        .post(`${this.$api_endpoints.FILES}${this.file.sha256}/actions/`, {
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
    if (this.$store.getters.getLikes.includes(this.file.sha256))
      this.liked = true
    if (this.file.timestamp)
      this.time = moment(this.file.timestamp).format("MMMM Do YYYY")
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
      #hash {
        font-size: large;
        font-weight: 500;
        cursor: pointer;
      }
      #Av {
        padding-left: 1em;
      }
      #timestamp {
        padding-left: 1em;
      }
      svg {
        vertical-align: bottom;
      }
      #tag {
        color: #00d1b2;
        font-weight: 600;
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
