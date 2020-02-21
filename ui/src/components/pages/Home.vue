<template>
  <div class="columns">
    <div class="column"></div>
    <div
      class="column is-8 box"
      v-if="usersData.length > 0 && activitiesToShow.length > 0"
    >
      <p id="no_activity" v-if="activities.length === 0">
        No Recent Activities
      </p>
      <ActivityCard
        :activity="activity"
        :userData="getUserDataPerActivity(activity.username)"
        v-for="(activity, index) in activitiesToShow"
        :key="index"
      />
    </div>
    <div class="column"></div>
  </div>
</template>

<script>
import ActivityCard from "../elements/home/ActivityCard"

export default {
  components: {
    ActivityCard,
  },
  data() {
    return {
      loged: false,
      activities: [],
      usersData: [],
      actToShowCount: 0,
    }
  },
  computed: {
    activitiesToShow: function() {
      return this.activities.slice(0, this.actToShowCount)
    },
  },
  methods: {
    getLoggedInActivities: function() {
      this.$http
        .get(
          this.$api_endpoints.USERS +
            this.$store.getters.getUsername +
            "/activities",
        )
        .then((res) => {
          this.formatActivities(res.data)
        })
        .catch(() => this.$awn.alert("Error Occured While getting activities"))
    },
    getLoggedOffActivities: function() {
      this.$http
        .get(this.$api_endpoints.USERS + "activities")
        .then((res) => this.formatActivities(res.data))
        .catch(() => this.$awn.alert("Error Occured While getting activities"))
    },
    formatActivities: function(activities) {
      for (var index in activities) {
        this.getUserData(activities[index].username)
        for (var index2 in activities[index].activities) {
          activities[index].activities[index2].username =
            activities[index].username
          this.activities.push(activities[index].activities[index2])
        }
      }
      this.orderActivities()
    },
    orderActivities: function() {
      this.activities.sort((a, b) => {
        return new Date(b.timestamp) - new Date(a.timestamp)
      })
    },
    getUserData: async function(username) {
      this.$http
        .get(this.$api_endpoints.USERS + username + "/avatar", {
          responseType: "arraybuffer",
        })
        .then((secRes) => {
          var data = {
            username: username,
            avatar: Buffer.from(secRes.data, "binary").toString("base64"),
          }
          this.usersData.push(data)
        })
        .catch(() => {
          this.$awn.alert("An Error Occured While fetshing the user data")
        })
    },
    getUserDataPerActivity(username) {
      return this._.find(this.usersData, (data) => data.username === username)
    },
    scroll() {
      window.onscroll = () => {
        var bottomOfWindow =
          document.documentElement.scrollTop + window.innerHeight ===
          document.documentElement.offsetHeight

        if (bottomOfWindow) {
          this.actToShowCount += 10
          if (this.actToShowCount > this.activities.length)
            this.actToShowCount = this.activities.length
        }
      }
    },
  },
  mounted() {
    if (this.$store.getters.getLoggedIn) {
      this.loged = true
      this.getLoggedInActivities()
    } else {
      this.getLoggedOffActivities()
      this.loged = false
    }
    this.actToShowCount = 10
    this.scroll()
  },
}
</script>

<style scoped>
.box {
  padding: 1.5em;
}
#no_activity {
  font-size: 20px;
  font-weight: 200;
}
</style>
