<template>
  <div>
    <form class="form">
      <h1 class="signin">Sign In</h1>
      <label for="email">Email</label>
      <input class="entry" id="email" type="email" v-model="email" autofocus placeholder="Email Address" />
      <label for="password">Password</label>
      <input class="entry" id="password" type="password" v-model="password" placeholder="Password" />
      <button class="login" type="submit" @click="handleSubmit">Sign In</button>

      <h3 class="forgot">
        <router-link  to="/forgot_password"
          >Forgot password?</router-link
        >
      </h3>
    </form>
    <h3 class="not-member">
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
  border-radius: 0.25rem;
  border: 1px solid #33333330;
}

.forgot,
.signin {
  margin: auto;
}


.signin {
  font-size: 2em;
  margin-bottom: 1em;
}

.not-member {
font-size: 16px;
text-align: center;
}


.form .entry {
  width: 280px;
  color: #333333;
  background: none;
  border: 1px solid #33333350;
  padding: 0.5em;
  font-size: inherit;
  border-radius: 0.25rem;
  box-shadow: inset 0 2px 4px 0 hsla(0, 0%, 0%, 0.08);
  margin-top: -1em;
}

.login {
  cursor: pointer;
  border-radius: 0.25rem;
  font-size: inherit;
  padding: 0.7em;
  font-weight: 600;
  color: white;
  background-color: #e7501d;
  border: none;
}
</style>
