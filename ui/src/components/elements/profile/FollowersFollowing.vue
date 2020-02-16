<template>
  <div v-if="active">
    <peopleCard
      v-for="(user, index) in usersData"
      :key="index"
      :userData="user"
    />
  </div>
</template>

<script>
import peopleCard from "./people"

export default {
  props: ["users", "active"],
  components: {
    peopleCard,
  },
  data() {
    return {
      usersData: [],
    }
  },
  watch: {
    users: function() {
      this.usersData = []
      for (var index in this.users) {
        this.getUserData(this.users[index])
      }
    },
  },
  methods: {
    getUserData: function(username) {
      if (username === this.$store.getters.getUsername) {
        var userData = this.$store.getters.getUserData
        var data = {
          username: username,
          name: userData.name,
          location: userData.location,
          avatar: this.$store.getters.getAvatar,
        }
        this.usersData.push(data)
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
              var data = {
                username: res.data.username,
                name: res.data.name,
                location: res.data.location,
                avatar: Buffer.from(secRes.data, "binary").toString("base64"),
              }
              this.usersData.push(data)
            })
        })
        .catch()
    },
  },
  mounted() {
    for (var index in this.users) {
      this.getUserData(this.users[index])
    }
  },
}
</script>

<style></style>
