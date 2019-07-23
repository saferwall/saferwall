<template>
  <div>
    <form class="form">
      <h1 class="signup">Create Your Account</h1>
      <div
        class="input-container"
        :class="{ valid: !$v.username.$invalid, 'not-valid': $v.username.$error }"
      >
        <label for="username">Username</label>

        <input
          class="entry"
          id="username"
          type="text"
          v-model="$v.username.$model"
          placeholder="Username"
          autocomplete="username"
        />
        <div v-show="$v.username.$dirty">
          <span v-show="!$v.username.required" class="error"
            >Username is required</span
          >
        </div>
      </div>
      <div
        class="input-container"
        :class="{
          valid: !$v.email.$invalid,
          'not-valid': $v.email.$error
        }"
      >
        <label for="email">Email</label>
        <input
          class="entry"
          id="email"
          type="email"
          v-model="$v.email.$model"
          placeholder="name@example.com"
          autocomplete="email"
        />
        <div v-show="$v.email.$dirty">
          <span v-show="!$v.email.required" class="error"
            >Email is required</span
          >

          <span v-show="!$v.email.email" class="error">Email is not valid</span>
          <!-- Add email in use constraint following backend implementation -->
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
          class="entry"
          id="password"
          type="password"
          v-model="$v.password.$model"
          placeholder="Password"
          autocomplete="new-password"
        />
        <div v-show="$v.password.$dirty">
          <span v-show="!$v.password.required" class="error"
            >Password is required</span
          >

          <span v-show="!$v.password.minLength" class="error"
            >Must be at least
            {{ $v.password.$params.minLength.min }} characters</span
          >
        </div>
      </div>
      <p>
        <input type="checkbox" id="tos" />
        <label for="tos"
          >&nbsp;I agree to the&nbsp;<router-link to="/tos"
            >Terms of Service</router-link
          ></label
        >
      </p>
      <button class="login" type="submit" @click="handleSubmit">
        Create Account
      </button>
    </form>
    <h3 class="already-member">
      Already have an account? <router-link to="/login">Sign in</router-link>
    </h3>
  </div>
</template>

<script>
import { required, minLength, email } from "vuelidate/lib/validators";

export default {
  data() {
    return {
      username: "",
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
    username: {
      required
    },
    email: {
      required,
      email
    },
    password: {
      required,
      minLength: minLength(8)
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

.signup {
  font-size: 2em;
  margin-bottom: 1em;
  text-align: center;
}

.already-member {
  font-size: 16px;
  text-align: center;
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
}

.form label {
  font-weight: 600;
}
.form label[for="tos"] {
  font-weight: normal;
  vertical-align: middle;
  font-size: 0.9em;
}
</style>
