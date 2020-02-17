<template>
  <div>
    <p id="no_user" v-if="!userData && !loading">No Such User Exists</p>
    <div class="columns tile is-ansestor box" v-if="userData">
      <div class="column is-one-quarter">
        <UserData :userData="userData" />
      </div>
      <div class="column">
        <div class="tabs is-medium ">
          <ul>
            <li
              :class="{ 'is-active': activeTab === 0 }"
              @click="activeTab = 0"
            >
              <a>Likes</a>
            </li>
            <li
              :class="{ 'is-active': activeTab === 1 }"
              @click="activeTab = 1"
            >
              <a>Submissions</a>
            </li>
            <li
              :class="{ 'is-active': activeTab === 2 }"
              @click="activeTab = 2"
            >
              <a>Followers</a>
            </li>
            <li
              :class="{ 'is-active': activeTab === 3 }"
              @click="activeTab = 3"
            >
              <a>Following</a>
            </li>
            <li
              :class="{ 'is-active': activeTab === 4 }"
              @click="activeTab = 4"
            >
              <a>Comments</a>
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
</style>
