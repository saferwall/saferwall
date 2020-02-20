<template>
  <div class="level tile">
    <div class="data">
      <p id="hash" @click="showFile">
        {{ this.hash }}
      </p>
      <p class="info">
        <span id="tags">
          <i class="icon fas fa-tags"></i>
          Tags:
          <span v-if="!fileData.tags">none</span>
          <span id="tag" v-for="tag in fileData.tags" :key="tag">{{
            tag
          }}</span>
        </span>
        <span id="Av">
          <i class="icon fas fa-search"></i>
          Av Detection Count: {{ fileData.AvDetectionCount }}
        </span>
      </p>
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
  props: ["hash"],
  data() {
    return {
      liked: false,
      fileData: {},
    }
  },
  watch: {
    hash: function() {
      if (this.$store.getters.getLikes.includes(this.hash)) this.liked = true
    },
  },
  methods: {
    likeUnlike: function() {
      this.$http
        .post(`${this.$api_endpoints.FILES}${this.hash}/actions/`, {
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
      this.$store.dispatch("updateHash", this.hash)
      this.$router.push(this.$routes.SUMMARY.path + this.hash)
    },
    getAvDetectionCount: function(scans) {
      var count = 0
      for (const av of Object.values(scans)) {
        if (av.infected) count++
      }
      return count
    },
    getFileData: function() {
      this.$http
        .get(
          this.$api_endpoints.FILES + this.hash + "?fields=sha256,tags,multiav",
        )
        .then((res) => {
          res.data.AvDetectionCount = this.getAvDetectionCount(
            res.data.multiav.last_scan,
          )
          this.fileData = res.data
        })
        .catch()
    },
  },
  mounted() {
    this.getFileData()
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
    font-weight: bold;
    cursor: pointer;
  }

  .info {
    svg {
      vertical-align: bottom;
    }
    #tag {
      color: #00d1b2;
      font-weight: 600;
    }
    #Av {
      padding-left: 0.5em;
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
