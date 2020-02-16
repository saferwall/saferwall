<template>
  <div>
    <article class="media box" v-if="!verification">
      <figure class="media-left">
        <p class="image is-64x64" @click="goToProfile">
          <img :src="'data:image/png;base64,' + data.avatar" />
        </p>
        <div class="info">
          <i class="icon fas fa-location-arrow"></i>
          <p>
            {{ this.data.location }}
          </p>
        </div>
      </figure>
      <div class="media-content">
        <div class="content">
          <p>
            <strong class="username" @click="goToProfile">{{
              this.data.name ? this.data.name : this.data.username
            }}</strong>
            &nbsp;
            <small>@{{ this.data.username }}</small>
            &nbsp;
            <small>{{ this.time }}</small>
            <br />
            <span class="comment_body" v-html="data.body"></span>
          </p>
        </div>
      </div>
      <div class="media-right">
        <button
          class="delete"
          v-if="deletable"
          @click="verification = true"
        ></button>
      </div>
    </article>
    <article class="media box" v-if="verification">
      <div class="media-content">
        <div class="content verif">
          <div>Are you Sure?</div>
          <button class="button is-light danger" @click="deleteComment">
            Yes
          </button>
          <button class="button is-light" @click="verification = false">
            No
          </button>
        </div>
      </div>
    </article>
  </div>
</template>

<script>
import moment from "moment"

export default {
  props: ["data"],
  data() {
    return {
      time: null,
      deletable: false,
      verification: false,
    }
  },
  methods: {
    formatTimestamp: function() {
      this.time = moment(this.data.timestamp).fromNow(true)
    },
    deleteComment: function() {
      this.$http
        .delete(
          this.$api_endpoints.FILES +
            this.$store.getters.getHashContext +
            "/comments/" +
            this.data.id,
        )
        .then(() => {
          this.$store.dispatch("updateComments")
        })
        .catch(() => {
          this.$awn.alert(
            "An Error Occured While Deleting the comment, try again",
          )
        })
    },
    goToProfile: function() {
      this.$router.push(this.$routes.PROFILE.path + this.data.username)
    },
  },
  mounted() {
    if (this.data.username !== this.$store.getters.getUsername) {
      this.deletable = false
    } else {
      this.deletable = true
    }
    this.formatTimestamp()
  },
}
</script>

<style scoped lang="scss">
.username {
  font-size: 1em;
  font-weight: 600;
  &:hover {
    color: #39d9c1;
    cursor: pointer;
  }
}

.image {
  &:hover {
    opacity: 0.5;
    cursor: pointer;
  }
}

.media {
  margin-bottom: 0.6em;

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
  .button {
    font-size: 1.5em;
    background-color: transparent;
    &:hover {
      background-color: #00000012;
    }
  }
  .danger {
    color: red;
    &:hover {
      color: red !important;
    }
  }
}

.info {
  display: flex;
  margin-top: 5px;
  p {
    padding-left: 3px;
    margin-right: auto;
  }
  svg {
    margin-left: auto;
  }
}
</style>
