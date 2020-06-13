<template>
  <div class="columns">
    <div class="column"></div>
    <div
      class="column is-8 box"
      v-if="usersData.length > 0 && activitiesToShow.length > 0"
    >
      <div class="title">
        <span>Activities</span>
      </div>
      <p id="no_activity" v-if="activities.length === 0">
        No Recent Activities
      </p>
      <ActivityCard
        :activity="activity"
        :userData="getUserDataPerActivity(activity.username)"
        :secondUser="
          activity.content.user
            ? getUserDataPerActivity(activity.content.user)
            : null
        "
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
      logged: false,
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
    getActivities: function() {
      const url = this.logged
        ? this.$api_endpoints.USERS +
          this.$store.getters.getUsername +
          "/activities"
        : this.$api_endpoints.USERS + "activities"
      this.$http
        .get(url)
        .then((res) => {
          // getting users data
          var users = this._.uniq(
            this._.flatten(
              this._.map(res.data, (activity) => {
                if (activity.content.user)
                  return [activity.username, activity.content.user]
                return [activity.username]
              }),
            ),
          )
          for (var user of users) {
            this.getUserData(user)
          }

          this.activities = res.data
          this.orderActivities()
        })
        .catch(() => this.$awn.alert("Error Occured While getting activities"))
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
      this.logged = true
    } else {
      this.logged = false
    }
    this.getActivities()
    this.actToShowCount = 10
    this.scroll()
  },
}
</script>

<style lang="scss" scoped>
.box {
  padding: 1.5em;
}
#no_activity {
  font-size: 20px;
  font-weight: 200;
}
.title {
  border-bottom: 1px solid;
  border-color: #ededed;
  padding-bottom: 0.6em;
  margin-bottom: 1.2em;
}
</style>
