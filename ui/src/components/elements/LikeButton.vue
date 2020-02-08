<template>
  <button
    class="button is-outlined is-danger"
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
  data() {
    return {
      liked: false,
    }
  },
  computed: {
    ...mapGetters({ likes: "getLikes" }),
  },
  watch: {
    hash: function(val){
      this.liked = this.likes.includes(val)
    }
  },
  methods: {
    likeUnlike: function() {
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
  mounted() {
    this.liked = this.likes.includes(this.hash)
  },
}
</script>

<style lang="scss" scoped>
.active {
  background-color: #f14668 !important;
  span{
    color: white;
  }
}
</style>
