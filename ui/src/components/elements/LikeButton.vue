<template>
  <button
    class="button is-outlined"
    :class="{ active: liked }"
    @click="likeUnlike"
  >
    <div class="rescan-button-text">
      <span class="icon">
        <i class="fas fa-heart"></i>
      </span>
      <span>
        {{ this.liked ? "Liked" : "Like" }}
      </span>
    </div>
  </button>
</template>

<script>
import { mapGetters } from "vuex"

export default {
  props: ["hash"],
  computed: {
    ...mapGetters({ likes: "getLikes" }),
    liked: function() {
      if (!this.likes) return false
      return this.likes.includes(this.hash)
    },
  },
  methods: {
    likeUnlike: function() {
      if (!this.$store.getters.getLoggedIn) {
        this.$router.push({
          name: "login",
          params: {
            nextUrl: this.$route.path,
          },
        })
        return
      }
      this.$http
        .post(`${this.$api_endpoints.FILES}${this.hash}/actions/`, {
          type: this.liked ? "unlike" : "like",
        })
        .then(() => {
          if (!this.liked) this.$store.dispatch("addRemoveLike", true)
          else this.$store.dispatch("addRemoveLike", false)
          this.liked = !this.liked
        })
        .catch(() => {
          this.$awn.alert("An Error Occured, try again")
        })
    },
  },
}
</script>

<style lang="scss" scoped>
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
</style>
