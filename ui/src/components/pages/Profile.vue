<template>
  <div>
    <p id="no_user" v-if="!userData && !loading">No Such User Exists</p>
    <p id="title" v-if="userData">User Profile:</p>
    <div class="columns" v-if="userData">
      <div class="column is-one-quarter box">
        <UserData :userData="userData" />
      </div>
      <div class="column box">
        <div class="tabs is-medium ">
          <ul>
            <li
              :class="{ 'is-active': activeTab === 0 }"
              @click="activeTab = 0"
            >
              <a
                >Likes
                <span class="counter">
                  {{ this.userData.likes ? this.userData.likes.length : "0" }}
                </span>
              </a>
            </li>
            <li
              :class="{ 'is-active': activeTab === 1 }"
              @click="activeTab = 1"
            >
              <a
                >Submissions
                <span class="counter">
                  {{
                    this.userData.submissions
                      ? this.userData.submissions.length
                      : "0"
                  }}
                </span>
              </a>
            </li>
            <li
              :class="{ 'is-active': activeTab === 2 }"
              @click="activeTab = 2"
            >
              <a
                >Followers
                <span class="counter">
                  {{
                    this.userData.followers
                      ? this.userData.followers.length
                      : "0"
                  }}
                </span>
              </a>
            </li>
            <li
              :class="{ 'is-active': activeTab === 3 }"
              @click="activeTab = 3"
            >
              <a
                >Following
                <span class="counter">
                  {{
                    this.userData.following
                      ? this.userData.following.length
                      : "0"
                  }}
                </span>
              </a>
            </li>
            <li
              :class="{ 'is-active': activeTab === 4 }"
              @click="activeTab = 4"
            >
              <a
                >Comments
                <span class="counter">
                  {{
                    this.userData.comments ? this.userData.comments.length : "0"
                  }}
                </span>
              </a>
            </li>
          </ul>
        </div>
        <Likes :active="activeTab === 0" :likes="userData.likes" />
        <Submissions :active="activeTab === 1" :subs="userData.submissions" />
        <Followers :active="activeTab === 2" :users="userData.followers" />
        <Following :active="activeTab === 3" :users="userData.following" />
        <Comments :active="activeTab === 4" :comments="userData.comments" />
      </div>
    </div>
  </div>
</template>

<script>
import UserData from "../elements/profile/UserData"
import Likes from "../elements/profile/Likes"
import Submissions from "../elements/profile/Submissions"
import Followers from "../elements/profile/Followers"
import Following from "../elements/profile/Following"
import Comments from "../elements/profile/Comments"

export default {
  components: {
    UserData,
    Likes,
    Followers,
    Following,
    Submissions,
    Comments,
  },
  data() {
    return {
      activeTab: 0,
      userData: null,
      loading: true,
    }
  },
  methods: {
    loadUseData: function(username) {
      if (username === this.$store.getters.getUsername) {
        this.userData = this.$store.getters.getUserData
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
              res.data.avatar = Buffer.from(secRes.data, "binary").toString(
                "base64",
              )
              this.userData = res.data
              this.loading = false
            })
        })
        .catch(() => {
          this.loading = false
          this.$awn.alert("An Error Occured While fetshing the user data")
        })
    },
  },
  mounted() {
    this.loadUseData(this.$route.params.user)
  },
  beforeRouteUpdate(to, from, next) {
    this.loading = true
    this.userData = null
    this.activeTab = 0
    this.loadUseData(to.params.user)
    next()
  },
}
</script>

<style scoped>
#no_user {
  font-size: 30px;
  font-weight: 200;
  text-align: center;
}
#title {
  font-size: 30px;
  font-weight: 200;
  margin-left: .5em;
  margin-bottom: .5em;
}
.counter {
  background-color: #f7f7f7;
  color: #4a4a5e;
  border-radius: 50%;
  min-width: 1.8em;
  height: 1.8em;
  text-align: center;
  margin-left: 5px;
  font-size: 0.8em;
  line-height: 180%;
}
.column {
  height: max-content;
  margin: 1em;
}
</style>
