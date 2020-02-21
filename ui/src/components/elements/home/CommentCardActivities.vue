<template>
  <div class="level tile">
    <div class="level-left">
      <div id="hash" @click="showFile">{{ data.sha256 }}</div>
      <div class="level-item">
        <span id="comment_body" v-html="data.body"></span>
      </div>
    </div>
    <div v-if="this.$store.getters.getLoggedIn">
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
</template>

<script>
export default {
  props: ["data"],
  data() {
    return {
      liked: false,
    }
  },
  watch: {
    data: function() {
      if (this.$store.getters.getLikes.includes(this.hash)) this.liked = true
    },
  },
  methods: {
    likeUnlike: function() {
      this.$http
        .post(`${this.$api_endpoints.FILES}${this.data.sha256}/actions/`, {
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
      this.$store.dispatch("updateHash", this.data.sha256)
      this.$router.push(this.$routes.SUMMARY.path + this.data.sha256)
    },
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
  padding: 0.7em;
  align-items: center;

  #hash {
    font-weight: bold !important;
    cursor: pointer;
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
