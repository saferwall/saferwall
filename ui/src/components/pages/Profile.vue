<template>
  <div class="columns tile is-ansestor box" v-if="userData">
    <div class="column is-one-quarter">
      <UserData :userData="userData" />
    </div>
    <div class="column">
      <div class="tabs is-medium ">
        <ul>
          <li :class="{ 'is-active': activeTab === 0 }" @click="activeTab = 0">
            <a>Likes</a>
          </li>
          <li :class="{ 'is-active': activeTab === 1 }" @click="activeTab = 1">
            <a>Submissions</a>
          </li>
          <li :class="{ 'is-active': activeTab === 2 }" @click="activeTab = 2">
            <a>Followers</a>
          </li>
          <li :class="{ 'is-active': activeTab === 3 }" @click="activeTab = 3">
            <a>Following</a>
          </li>
          <li :class="{ 'is-active': activeTab === 4 }" @click="activeTab = 4">
            <a>Comments</a>
          </li>
        </ul>
      </div>
      <Likes v-if="activeTab === 0" :likes="userData.likes" />
    </div>
  </div>
</template>

<script>
import UserData from "../elements/profile/UserData"
import Likes from "../elements/profile/Likes"

export default {
  components: {
    UserData,
    Likes,
  },
  data() {
    return {
      activeTab: 0,
      userData: null,
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
            })
        })
        .catch(() => {
          this.$awn.alert("An Error Occured While fetshing the user data")
        })
    },
  },
  mounted() {
    this.loadUseData(this.$route.params.user)
  },
}
</script>

<style></style>
