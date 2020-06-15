<template>
  <div class="level tile">
    <div class="level-left">
      <div id="hash" @click="showFile">{{ data.content.sha256 }}</div>
      <div class="level-item">
        <span id="comment_body" v-html="data.content.body"></span>
      </div>
      <p class="info">
        <span id="tags">
          <i class="icon fas fa-tags"></i>
          Tags:
          <span v-if="!data.tags">none</span>
          <span
            v-else
            class="tag is-link is-normal"
            :class="{ redTag: isAntivirusTag(tag[0]) }"
            id="tag"
            v-for="tag in tags"
            :key="tag[1]"
          >
            {{ tag[1] }}
          </span>
        </span>
        <span id="Av">
          <i class="icon fas fa-shield-alt"></i>
          Antivirus: {{ data.av_count }}/12
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
  props: ["data"],
  data() {
    return {
      liked: false,
    }
  },
  watch: {
    data: function() {
      if (this.$store.getters.getLikes.includes(this.data.content.sha256))
        this.liked = true
    },
  },
  computed: {
    tags: function() {
      var tags = []
      if (!this.data.tags) return null
      for (var tag of Object.entries(this.data.tags)) {
        for (var value of tag[1]) {
          tags.push([tag[0], value.toLowerCase()])
        }
      }
      return this._.uniqWith(tags, (x, y) => x[1] === y[1])
    },
  },
  methods: {
    likeUnlike: function() {
      this.$http
        .post(
          `${this.$api_endpoints.FILES}${this.data.content.sha256}/actions/`,
          {
            type: this.liked ? "unlike" : "like",
          },
        )
        .then(() => {
          this.liked = !this.liked
          this.$store.dispatch("updateLikes")
        })
        .catch(() => {
          this.$awn.alert("An Error Occured, try again")
        })
    },
    showFile: function() {
      this.$store.dispatch("updateHash", this.data.content.sha256)
      this.$router.push(this.$routes.SUMMARY.path + this.data.content.sha256)
    },
    isAntivirusTag: function(tag) {
      const antivirusList = [
        "eset",
        "fsecure",
        "avira",
        "bitdefender",
        "kaspersky",
        "symantec",
        "sophos",
        "windefender",
        "clamav",
        "comodo",
        "avast",
        "mcafee",
      ]
      return antivirusList.includes(tag)
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
  .info {
    margin-top: 0.5rem;
    svg {
      vertical-align: bottom;
    }
    #tag {
      margin-right: 0.2em;
      color: white;
      font-weight: 600;
    }
    .redTag {
      background-color: #f14668;
    }
    #Av {
      padding-left: 0.5em;
    }
  }
}
</style>
