<template>
  <div>
    <form class="form">
      <h1 class="signin">Sign In</h1>
      <label for="username">Username </label>
      <input id="username" type="text" v-model="username" />
      <label for="password">Password </label>
      <input id="password" type="password" v-model="password" />
      <button type="submit" @click="handleSubmit">Sign In</button>

      <h3 class="forgot">
        <router-link to="/forgot_password">Forgot password?</router-link>
      </h3>
      
    </form><h3 class="signup">
        Not a member? <router-link to="/signup">Sign up</router-link>
      </h3>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: "",
      password: ""
    };
  },
  methods: {
    handleSubmit(e) {
      e.preventDefault();

      if (this.password === "user" && this.username === "user") {
        this.$router.push(this.$route.params.nextUrl);
      }
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.form {
  display: grid;
  line-height: 2em;
  align-items: center; /* align-self every label item vertically in its row!*/
  justify-content: center;
  width: max-content;
  grid-row-gap: 1em;
  padding: 4em;
  color: #333333;
  background-color: white;
  font-size: 16px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1),
    0 4px 6px -2px rgba(0, 0, 0, 0.05);
}

.forgot,
.signin {
  margin: auto;
}

.signin {
  font-size: 2em;
  margin-bottom: 1em;
}

.signup {
}

.form input,
select {
  width: 280px;
  color: #d6d6d638;
  background: none;
  padding: 0.5em;
  font-size: inherit;
}

button {
  font-size: inherit;
  padding: 0.5em;
}
</style>
