<template>
  <div>
    <form novalidate="true" class="form" @submit="handleSubmit">
      <h1 class="signin">Sign In</h1>
      <div
        class="input-container"
        :class="{ valid: !$v.email.$invalid, 'not-valid': $v.email.$error }"
      >
        <label for="email">Email</label>
        <input
          required
          class="entry"
          id="email"
          type="email"
          v-model.trim="$v.email.$model"
          placeholder="name@example.com"
          autocomplete="email"
        />
        <div v-show="$v.email.$dirty">
          <span v-show="!$v.email.required" class="error"
            >Email is required</span
          >
          <span v-show="!$v.email.email" class="error">Email is not valid</span>
        </div>
      </div>
      <div
        class="input-container"
        :class="{
          valid: !$v.password.$invalid,
          'not-valid': $v.password.$error
        }"
      >
        <label for="password">Password</label>
        <input
          required
          class="entry"
          id="password"
          type="password"
          v-model.trim="$v.password.$model"
          placeholder="Password"
          autocomplete="current-password"
        />

        <div v-show="$v.password.$dirty">
          <span v-show="!$v.password.required" class="error"
            >Password is required</span
          >
        </div>
      </div>
      <button class="login" type="submit">Sign In</button>
      <h3 class="forgot">
        <router-link to="/forgot_password">Forgot password?</router-link>
      </h3>
    </form>
    <h3 class="not-member">
      Not a member? <router-link to="/signup">Sign up</router-link>
    </h3>
  </div>
</template>

<script>
import { required, email } from "vuelidate/lib/validators";

export default {
  data() {
    return {
      email: "",
      password: ""
    };
  },
  methods: {
    handleSubmit(e) {
      e.preventDefault();
    }
  },
  validations: {
    email: {
      required,
      email
    },
    password: {
      required,
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
  grid-row-gap: .1em;
  padding: 4em;
  color: #333333;
  background-color: white;
  font-size: 16px;
  border-radius: 0.25rem;
  border: 1px solid #33333330;
}

.valid input {
  background: url("../../assets/imgs/check.svg") !important;
  background-repeat: no-repeat !important;
  background-position: 97% center !important;
  border: 1px solid #02a678 !important;
}

.not-valid input {
  background: url("../../assets/imgs/error.svg") !important;
  background-repeat: no-repeat !important;
  background-position: 97% center !important;
  border: 1px solid #bb0000 !important;
}

.input-container {
  display: flex;
  height: 100px;
  flex-direction: column;
}

.error {
  font-size: 0.8em;
  color: #bb0000;
  font-weight: 100;
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
  width: 340px;
  min-height: 45px;
  color: #333333;
  background: none;
  border: 1px solid #33333335;
  padding: 0.5em;
  font-size: inherit;
  border-radius: 0.25rem;
  box-shadow: inset 0 2px 4px 0 hsla(0, 0%, 0%, 0.03);
  transition: border 0.1s ease;
}

.form .entry:hover {
  border: 1px solid #3333336b;
}

.form .entry:focus {
  outline: none;
  border: 1px solid #3333336b;
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
  margin-top: 1em;
}

.form label {
  font-weight: 600;
}
</style>
