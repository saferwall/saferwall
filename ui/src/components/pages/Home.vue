<template>
  <div class="columns">
    <div class="column"></div>
    <div class="column is-8 box">
      <ActivityCard
        :activity="activity"
        :userData="getUserDataPerActivity(activity.username)"
        v-for="(activity, index) in activities"
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
    }
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
        .get(this.$api_endpoints.USERS + username)
        .then((res) => {
          this.$http
            .get(this.$api_endpoints.USERS + username + "/avatar", {
              responseType: "arraybuffer",
            })
            .then((secRes) => {
              var data = {
                username: res.data.username,
                name: res.data.name,
                location: res.data.location,
                avatar: Buffer.from(secRes.data, "binary").toString("base64"),
              }
              this.usersData.push(data)
            })
        })
        .catch(() => {
          this.$awn.alert("An Error Occured While fetshing the user data")
        })
    },
    getUserDataPerActivity(username) {
      return this._.find(this.usersData, (data) => data.username === username)
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
  },
}
</script>

<style scoped></style>
